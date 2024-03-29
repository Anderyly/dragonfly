/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package api

import (
	"dragonfly/ay"
	"dragonfly/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat"
	"log"
	"net/http"
	"time"
)

type NotifyController struct {
}

func (con NotifyController) Wechat(c *gin.Context) {

	notifyReq, err := wechat.ParseNotifyToBodyMap(c.Request)
	con.CheckErr(err)
	log.Println(notifyReq)

	j, err := json.Marshal(notifyReq)
	con.CheckErr(err)

	ay.CreateMutiDir("log/wechat")
	ay.WriteFile("log/wechat/"+notifyReq.Get("out_trade_no")+".txt", []byte(j), 0755)

	// 验签操作
	ok, err := wechat.VerifySign(ay.Yaml.GetString("wechat.key"), wechat.SignType_MD5, notifyReq)

	if !ok {
		log.Println(err)
		c.String(http.StatusOK, "%s", "fail")
		return
	}

	var order models.Order
	ay.Db.First(&order, "out_trade_no = ?", notifyReq.Get("out_trade_no"))

	// 查询订单失败
	if order.Id == 0 {
		rsp := new(wechat.NotifyResponse) // 回复微信的数据
		rsp.ReturnCode = gopay.FAIL
		rsp.ReturnMsg = gopay.FAIL
		c.String(http.StatusOK, "%s", rsp.ToXmlString())
		return
	}

	// 订单已处理过
	if order.Status == 1 {
		rsp := new(wechat.NotifyResponse) // 回复微信的数据
		rsp.ReturnCode = gopay.SUCCESS
		rsp.ReturnMsg = gopay.OK
		c.String(http.StatusOK, "%s", rsp.ToXmlString())
		return
	}

	res := 0
	row := false

	switch order.Type {
	case 5:
		// 增加余额
		res = con.AdduserAmount(order)
	case 4:
		// 训练营
		if err := ay.Db.Create(&models.TrainLog{
			TrainId: order.Cid,
			Uid:     order.Uid,
			Oid:     order.Id,
		}).Error; err != nil {
			res = 0
		} else {
			var user models.User
			ay.Db.First(&user, order.Uid)
			isVip, _ := CommonController{}.GetVip(user)
			if isVip {
				user.VipNum += order.Amount
				ay.Db.Save(&user)
			}
			res = 1
		}
	case 3:
		// 会员
		var user models.User
		ay.Db.First(&user, order.Uid)
		var vip models.Vip
		ay.Db.First(&vip, order.Cid)
		row = models.UserModel{}.SetVip(user, vip.Effective, vip.Amount)
		var user1 models.User
		ay.Db.First(&user1, order.Uid)
		CommonController{}.Give(user1)
	case 2:
		var card models.Card
		ay.Db.First(&card, order.Cid)
		row = models.UserCardModel{}.Add(card.Id, card.AllHour, order.Uid, card.Effective, order.OldAmount)
	case 1:
		// 舞蹈室
		type cv struct {
			Num  int      `json:"num"`
			Time []string `json:"time"`
			Ymd  int      `json:"ymd"`
		}
		var cc cv
		json.Unmarshal([]byte(order.Content), &cc)

		if order.CardId != 0 {
			// 扣除次卡
			var card models.UserCard
			ay.Db.First(&card, order.CardId)
			card.UseHour += cc.Num
			if card.UseHour >= card.AllHour {
				card.Status = 1
			}
			ay.Db.Save(&card)
		}
		// 增加预约次数
		isSub, _ := models.UserRoomSubscribeModel{}.Add(order.OutTradeNo, order.Content, order.Uid, order.Cid, cc.Ymd)
		if !isSub {
			ay.Json{}.Msg(c, 400, "请联系管理员", gin.H{})
			return
		}

		var user models.User
		ay.Db.First(&user, order.Uid)
		isVip, _ := CommonController{}.GetVip(user)
		if isVip {
			user.VipNum += order.Amount
			ay.Db.Save(&user)
		}
		res = 1

	}

	if res == 1 || row == true {
		stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"), time.Local)
		// 优惠卷设置过期
		var coupon models.Coupon
		ay.Db.First(&coupon, "id = ?", order.CouponId)
		if coupon.Id != 0 {
			coupon.Status = 1
			coupon.UseAt = models.MyTime{Time: stamp}
			ay.Db.Save(&coupon)
		}

		order.PayType = 1
		order.Status = 1
		order.TradeNo = notifyReq.Get("transaction_id")
		order.PayAt = models.MyTime{Time: stamp}
		ay.Db.Save(&order)

		rsp := new(wechat.NotifyResponse) // 回复微信的数据
		rsp.ReturnCode = gopay.SUCCESS
		rsp.ReturnMsg = gopay.OK
		c.String(http.StatusOK, "%s", rsp.ToXmlString())
	} else {
		rsp := new(wechat.NotifyResponse) // 回复微信的数据
		rsp.ReturnCode = gopay.FAIL
		rsp.ReturnMsg = gopay.FAIL
		c.String(http.StatusOK, "%s", rsp.ToXmlString())
	}
}

func (con NotifyController) CheckErr(err error) {
	log.Println(err)
}

func (con NotifyController) AdduserAmount(order models.Order) int {

	var user models.User
	ay.Db.First(&user, order.Uid)
	if user.Id == 0 {
		return 0
	}

	amount := user.Amount + order.OldAmount

	var rechargeAmount models.Amount
	ay.Db.Where("amount = ?", fmt.Sprintf("%.2f", order.OldAmount)).First(&rechargeAmount)
	if rechargeAmount.Id != 0 {
		// 赠送余额
		amount += rechargeAmount.GiveAmount
	}

	user.Amount = amount
	//user.VipNum += order.OldAmount
	tx := ay.Db.Begin()
	log.Println(user)
	if err := tx.Save(&user).Error; err != nil {
		log.Println(err)
		tx.Rollback()
		return 0
	} else {
		tx.Commit()
		log.Println("success")
		return 1
	}
}
