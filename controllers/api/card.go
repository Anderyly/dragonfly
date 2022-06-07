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
	"github.com/gin-gonic/gin"
	"time"
)

type CardController struct {
}

type listCardControllerForm struct {
	Page int `form:"page" binding:"required" label:"页码"`
}

func (con CardController) List(c *gin.Context) {
	var data listCardControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	page := data.Page - 1

	var coupon []models.Card
	ay.Db.Offset(page * Limit).
		Limit(Limit).
		Order("created_at asc").
		Find(&coupon)

	var list []gin.H
	for _, v := range coupon {
		list = append(list, gin.H{
			"id":             v.Id,
			"type":           v.Type,
			"name":           v.Name,
			"amount":         v.Amount,
			"hour_amount":    v.HourAmount,
			"all_hour":       v.AllHour,
			"effective_text": v.EffectiveText,
		})
	}

	if list == nil {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"list": []string{},
		})
	} else {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"list": list,
		})
	}
}

type userCardControllerForm struct {
	Page int `form:"page" binding:"required" label:"页码"`
}

func (con CardController) User(c *gin.Context) {
	var data userCardControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	isAuth, code, msg, requestUser := Auth(GetToken(Token))

	if !isAuth {
		ay.Json{}.Msg(c, code, msg, gin.H{})
		return
	}

	page := data.Page - 1

	type r struct {
		Id          int64         `json:"id"`
		Type        int           `json:"type"`
		Name        string        `json:"name"`
		Amount      float64       `json:"amount"`
		HourAmount  float64       `json:"hour_amount"`
		AllHour     int           `json:"all_hour"`
		EffectiveAt models.MyTime `json:"effective_at"`
		UseHour     int           `json:"use_hour"`
		Status      int           `json:"status"`
	}
	var res []r
	ay.Db.Table("d_user_card").
		Select("d_card.type,d_card.name,d_card.amount,d_card.hour_amount,d_card.all_hour,d_user_card.use_hour,d_user_card.status,d_user_card.effective_at").
		Joins("left join d_card on d_user_card.card_id=d_card.id").
		Offset(page*Limit).
		Limit(Limit).
		Where("d_user_card.uid = ?", requestUser.Id).
		Order("d_user_card.status asc").
		Find(&res)

	var list []gin.H
	for _, v := range res {
		hour := v.AllHour - v.UseHour
		status := v.Status
		if hour <= 0 {
			v.Status = 1
			status = 1
			ay.Db.Save(&v)
		}
		if v.EffectiveAt.Unix() < time.Now().Unix() {
			v.Status = 2
			status = 2
			ay.Db.Save(&v)
		}
		list = append(list, gin.H{
			"id":           v.Id,
			"type":         v.Type,
			"name":         v.Name,
			"amount":       v.Amount,
			"hour_amount":  v.HourAmount,
			"all_hour":     hour,
			"effective_at": v.EffectiveAt.Unix(),
			"status":       status,
		})
	}

	if list == nil {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"list": []string{},
		})
	} else {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"list": list,
		})
	}
}

type payCardControllerForm struct {
	Id       int64  `form:"id" binding:"required" label:"次卡id"`
	CouponId int64  `form:"coupon_id" label:"优惠卷id"`
	Code     string `form:"code" binding:"required" label:"微信code"`
	PayType  int    `form:"pay_type" binding:"required" label:"支付方式"`
}

func (con CardController) Pay(c *gin.Context) {
	var data payCardControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	isAuth, code, msg, requestUser := Auth(GetToken(Token))

	if !isAuth {
		ay.Json{}.Msg(c, code, msg, gin.H{})
		return
	}

	// 获取会员卡信息
	var card models.Card
	ay.Db.First(&card, data.Id)
	if card.Id == 0 {
		ay.Json{}.Msg(c, 400, "次卡不存在", gin.H{})
		return
	}

	// 优惠卷判断
	isCoupon, couponMsg, coupon := CommonController{}.IsCoupon(requestUser.Id, data.CouponId, 5)

	if data.CouponId != 0 {
		if !isCoupon {
			ay.Json{}.Msg(c, 400, couponMsg, gin.H{})
			return
		}
	}

	amount := card.Amount - coupon.Amount

	isOrder, outTradeNo, order := CommonController{}.MakeOrder(2, 2, "购买次卡", requestUser.Id, coupon.Id, amount, card.Amount, c.ClientIP(), "", 0, card.Id, 0)
	if isOrder == 0 {
		ay.Json{}.Msg(c, 400, outTradeNo, gin.H{})
		return
	}

	if data.PayType == 1 {
		isC, msg, info := PayController{}.Wechat(data.Code, order.Des, order.OutTradeNo, c.ClientIP(), amount)

		if !isC {
			ay.Json{}.Msg(c, 400, msg, gin.H{})
			return
		}

		ay.Json{}.Msg(c, 200, "success", gin.H{
			"info": info,
		})
	} else {
		is, msg, _ := CommonController{}.UserAmountPay(&requestUser, amount)
		if is {
			tx := ay.Db.Begin()
			order.Status = 1
			stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"), time.Local)
			order.PayAt = models.MyTime{Time: stamp}
			if err := tx.Save(&order).Error; err != nil {
				tx.Rollback()
				ay.Json{}.Msg(c, 400, "请联系管理员", gin.H{})
				return
			}
			isUserCard := models.UserCardModel{}.Add(card.Id, card.AllHour, requestUser.Id, card.Effective, card.Amount)
			if !isUserCard {
				tx.Rollback()
				ay.Json{}.Msg(c, 400, "请联系管理员", gin.H{})
				return
			}
			tx.Commit()
			ay.Json{}.Msg(c, 200, "用户支付成功", gin.H{})
		} else {
			ay.Json{}.Msg(c, 400, msg, gin.H{})

		}
	}

}

type selectCardControllerForm struct {
	Page int `form:"page" binding:"required" label:"页码"`
	Type int `form:"type" binding:"required" label:"类型"`
}

func (con CardController) Select(c *gin.Context) {
	var data selectCardControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	isAuth, code, msg, requestUser := Auth(GetToken(Token))

	if !isAuth {
		ay.Json{}.Msg(c, code, msg, gin.H{})
		return
	}

	page := data.Page - 1

	type r struct {
		Id          int64         `json:"id"`
		Type        int           `json:"type"`
		Name        string        `json:"name"`
		Amount      float64       `json:"amount"`
		HourAmount  float64       `json:"hour_amount"`
		AllHour     int           `json:"all_hour"`
		EffectiveAt models.MyTime `json:"effective_at"`
		UseHour     int           `json:"use_hour"`
		Status      int           `json:"status"`
	}
	var res []r
	ay.Db.Table("d_user_card").
		Select("d_card.type,d_card.name,d_card.amount,d_card.hour_amount,d_card.all_hour,d_user_card.use_hour,d_user_card.status,d_user_card.effective_at").
		Joins("left join d_card on d_user_card.card_id=d_card.id").
		Offset(page*Limit).
		Limit(Limit).
		Where("d_user_card.uid = ? and d_card.type = ?", requestUser.Id, data.Type).
		Order("d_user_card.status asc").
		Find(&res)

	var list []gin.H
	for _, v := range res {
		hour := v.AllHour - v.UseHour
		status := v.Status
		if hour <= 0 {
			v.Status = 1
			status = 1
			ay.Db.Save(&v)
		}
		if v.EffectiveAt.Unix() < time.Now().Unix() {
			v.Status = 2
			status = 2
			ay.Db.Save(&v)
		}
		list = append(list, gin.H{
			"id":           v.Id,
			"type":         v.Type,
			"name":         v.Name,
			"amount":       v.Amount,
			"hour_amount":  v.HourAmount,
			"all_hour":     hour,
			"effective_at": v.EffectiveAt.Unix(),
			"status":       status,
		})
	}

	if list == nil {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"list": []string{},
		})
	} else {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"list": list,
		})
	}
}
