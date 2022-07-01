/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package api

import (
	"context"
	"dragonfly/ay"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/pkg/util"
	"github.com/go-pay/gopay/wechat"
	"strconv"
	"strings"
	"time"
)

type PayController struct {
}

func (con PayController) Wechat(code, des, outTradeNo, ip string, amount float64) (bool, string, gin.H) {
	isCode, openid := LoginController{}.GetOpenid(code)
	if !isCode {
		return false, openid, gin.H{}
	}

	ctx, _ := context.WithCancel(context.Background())

	client := wechat.NewClient(ay.Yaml.GetString("wechat.appid"), ay.Yaml.GetString("wechat.mchId"), ay.Yaml.GetString("wechat.key"), true)
	// 打开Debug开关，输出请求日志，默认关闭
	client.DebugSwitch = gopay.DebugOn
	client.SetCountry(wechat.China)

	tmpStr1 := fmt.Sprintf("%.2f", amount)
	total, _ := strconv.ParseInt(strings.Replace(tmpStr1, ".", "", 1), 10, 64)
	//total := int(tmpnum1 * 100)
	//log.Println(tmpnum1)
	//log.Println(total)

	bm := make(gopay.BodyMap)
	bm.Set("nonce_str", util.RandomString(32)).
		Set("body", des).
		Set("out_trade_no", outTradeNo).
		Set("total_fee", total).
		Set("spbill_create_ip", ip).
		Set("notify_url", ay.Yaml.GetString("domain")+"/api/notify/wechat").
		Set("trade_type", "JSAPI").
		Set("sign_type", "MD5").
		Set("openid", openid)

	wxRsp, err := client.UnifiedOrder(ctx, bm)
	if err != nil {
		return false, "请联系管理员", gin.H{}
	}

	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	packages := "prepay_id=" + wxRsp.PrepayId // 此处的 wxRsp.PrepayId ,统一下单成功后得到
	paySign := wechat.GetH5PaySign(ay.Yaml.GetString("wechat.appid"), wxRsp.NonceStr, packages, wechat.SignType_MD5, timeStamp, ay.Yaml.GetString("wechat.key"))

	return true, "success", gin.H{
		"appId":     wxRsp.Appid,
		"timeStamp": timeStamp,
		"nonceStr":  wxRsp.NonceStr,
		"package":   packages,
		"signType":  "MD5",
		"paySign":   paySign,
	}
}

func (con PayController) WechatRefund(outTradeNo string, amount float64) (bool, string) {
	ctx, _ := context.WithCancel(context.Background())

	client := wechat.NewClient(ay.Yaml.GetString("wechat.appid"), ay.Yaml.GetString("wechat.mchId"), ay.Yaml.GetString("wechat.key"), true)
	// 打开Debug开关，输出请求日志，默认关闭
	client.DebugSwitch = gopay.DebugOn
	client.SetCountry(wechat.China)

	// 添加微信pem证书
	client.AddCertPemFilePath("cert/apiclient_cert.pem", "cert/apiclient_key.pem")
	//client.AddCertPemFileContent()

	bm := make(gopay.BodyMap)
	bm.Set("nonce_str", util.RandomString(32)).
		Set("out_trade_no", outTradeNo).
		Set("out_refund_no", outTradeNo).
		Set("total_fee", amount*100).
		Set("refund_fee", amount*100).
		Set("notify_url", ay.Yaml.GetString("domain")+"/api/notify/refund")

	_, _, err := client.Refund(ctx, bm)
	if err != nil {
		//return false, "请联系管理员", gin.H{}
		return false, err.Error()
		//log.Println(err)
	} else {
		return true, "success"
	}
	//log.Println(wxRsp)
	//log.Println(resBm)
}
