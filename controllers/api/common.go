package api

import (
	"dragonfly/ay"
	"dragonfly/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	qrcode "github.com/skip2/go-qrcode"
	"strconv"
	"strings"
	"time"
)

var (
	Token string
	Limit = 10
)

type CommonController struct {
}

func (con CommonController) A() {

}

func GetRequestIP(c *gin.Context) string {
	reqIP := c.ClientIP()
	if reqIP == "::1" {
		reqIP = "127.0.0.1"
	}
	return reqIP
}

func GetToken(token string) string {
	uid := ay.AuthCode(token, "DECODE", "", 0)
	return uid
}

func Auth(uid string) (bool, int, string, models.User) {
	var user models.User
	ay.Db.First(&user, "id = ?", uid)

	if user.Id == 0 {
		return false, 401, "token错误", user
	}
	return true, 200, "success", user

}

func (con CommonController) MakeQrCode(text string) string {
	name := ay.MD5(fmt.Sprintf("%s%s", text, time.Now().String())) + ".png"
	fileDir := fmt.Sprintf("static/qrcode/%d-%d/", time.Now().Year(), time.Now().Month())

	err := ay.CreateMutiDir(fileDir)
	if err != nil {
		return ""
	}

	err = qrcode.WriteFile(text, qrcode.Medium, 152, fileDir+name)
	if err != nil {
		return ""
	}
	return fileDir + name

}

// IsCoupon 验证优惠卷
func (con CommonController) IsCoupon(uid int64, couponId int64, productId int) (bool, string, models.Coupon) {

	var coupon models.Coupon
	ay.Db.First(&coupon, couponId)
	if coupon.Id == 0 {
		return false, "优惠卷不存在", coupon
	}

	if coupon.Uid != uid {
		return false, "不是您的优惠卷", coupon
	}

	if coupon.Status != 0 {
		return false, "优惠卷已使用", coupon
	}

	if coupon.EffectiveAt.Unix() < time.Now().Unix() {
		return false, "优惠卷已过期", coupon
	}

	is := 0
	for _, v := range strings.Split(coupon.Type, ",") {
		if v == strconv.Itoa(productId) {
			is = 1
			break
		}
	}

	if is == 0 {
		return false, "优惠卷不适用于此产品", coupon
	}

	return true, "success", coupon
}

func (con CommonController) GetVip(user models.User) (bool, models.VipLevel) {

	isVip := false

	if user.EffectiveAt.Unix() > time.Now().Unix() {
		isVip = true
	}

	var vipLevel models.VipLevel

	ay.Db.Order("num desc").
		Where("num <= ?", user.VipNum).
		First(&vipLevel)

	return isVip, vipLevel
}

// MakeOrder 生成订单号
func (con CommonController) MakeOrder(op, vType int, msg string, uid, CouponId int64, amount, OldAmount float64, ip, content string, cardId int64, cid int64, vipDiscount float64) (int, string, models.Order) {
	oid := ay.MakeOrder(time.Now())

	order := models.Order{
		OutTradeNo:  oid,
		Type:        vType,
		Ip:          ip,
		Des:         msg,
		Amount:      amount,
		Uid:         uid,
		Status:      0,
		PayType:     2,
		OldAmount:   OldAmount,
		CouponId:    CouponId,
		Content:     content,
		CardId:      cardId,
		Cid:         cid,
		Op:          op,
		VipDiscount: vipDiscount,
	}

	ay.Db.Create(&order)

	if order.Id == 0 {
		return 0, "数据错误，请联系管理员", order
	} else {
		return 1, oid, order
	}
}

func (con CommonController) GetChanelCount(channelNum int, uid int64) (bool, string) {
	// 查询当月取消次数
	var channelCount int64
	ay.Db.Model(models.UserChannel{}).
		Where("(DATE_FORMAT(created_at, '%Y%m') = ?) AND uid = ?", time.Now().Format("200601"), uid).
		Count(&channelCount)

	if channelNum <= int(channelCount) {
		return false, "当前取消已上限"
	} else {
		return true, "success"
	}
}

func (con CommonController) Back(oid int64) (bool, string) {

	var order models.Order
	ay.Db.First(&order, oid)

	if order.Status != 1 {
		return false, "订单未付款"
	}

	var user models.User
	ay.Db.First(&user, order.Uid)

	tx := ay.Db.Begin()

	if order.Type != 5 {
		user.VipNum -= order.Amount
	}

	// 退还金额 余额
	if order.PayType == 2 {
		user.Amount += order.Amount
	} else {
		is, msg := PayController{}.WechatRefund(order.OutTradeNo, order.Amount)
		if !is {
			tx.Rollback()
			return false, msg
		}
	}

	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return false, "退还失败"
	}

	// 退还优惠卷
	if order.CouponId != 0 {
		var coupon models.Coupon
		ay.Db.First(&coupon, order.CouponId)
		coupon.Status = 0
		coupon.UseAt = models.MyTime{}
		if err := tx.Save(&coupon).Error; err != nil {
			tx.Rollback()
			return false, "优惠卷退还失败"
		}
	}

	// 退还次卡
	if order.CardId != 0 {
		type cv struct {
			Num  int      `json:"num"`
			Time []string `json:"time"`
			Ymd  int      `json:"ymd"`
		}
		var cc cv
		json.Unmarshal([]byte(order.Content), &cc)

		var card models.UserCard
		ay.Db.Where("id = ?", order.CardId).First(&card)
		card.Status = 0
		card.UseHour -= cc.Num
		if err := tx.Save(&card).Error; err != nil {
			tx.Rollback()
			return false, "次卡退还失败"
		}
	}

	order.Status = 2
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return false, "订单处理失败"
	}
	tx.Commit()

	return true, "success"
}

func (con CommonController) UserAmountPay(user *models.User, amount float64) (bool, string, *models.User) {

	if user.Amount < amount {
		return false, "用户余额不足", user
	}

	user.Amount -= amount
	if err := ay.Db.Save(&user).Error; err != nil {
		return false, "请联系管理员", user
	} else {
		return true, "success", user
	}

}

// Give 发放会员卡以及优惠卷
func (con CommonController) Give(user models.User) {

	// 获取当前时间
	ymd := time.Now().Unix()

	// 查询最后发放时间
	var lastGive models.UserVipGive
	ay.Db.Where("uid = ?", user.Id).Order("created_at desc").First(&lastGive)
	if lastGive.CreatedAt.Unix()+30*24*3600 < ymd {

		_, vipLevel := con.GetVip(user)
		// 获取系统配置
		config := models.ConfigModel{}.Get()

		// 获取会员卡
		var card models.Card
		ay.Db.First(&card, vipLevel.GiveCardId)
		if card.Id == 0 {
			return
		}

		// 发放次卡
		models.UserCardModel{}.Add(card.Id, card.AllHour, user.Id, card.Effective, 0)

		// 发放优惠卷
		for i := 0; i < vipLevel.GiveCoupon; i++ {
			effective := time.Now().Unix() + int64(config.CouponEffective*3600*24)
			stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Unix(effective, 0).Format("2006-01-02 15:04:05"), time.Local)
			ay.Db.Create(&models.Coupon{
				Type:        config.CouponType,
				Uid:         user.Id,
				Name:        config.CouponName,
				SubName:     config.CouponSubName,
				Des:         config.CouponDes,
				Amount:      config.CouponAmount,
				EffectiveAt: models.MyTime{Time: stamp},
				Status:      0,
			})
		}

		y, _ := strconv.Atoi(time.Now().Format("20060102"))
		ay.Db.Create(&models.UserVipGive{
			Uid: user.Id,
			Ymd: y,
		})
	}
}
