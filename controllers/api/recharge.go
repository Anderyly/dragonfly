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
	"strconv"
)

type RechargeController struct {
}

func (con RechargeController) Amount(c *gin.Context) {

	var res []models.Amount

	ay.Db.Order("sort asc").Find(&res)

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"list": res,
	})
}

type payRechargeControllerForm struct {
	Amount   float64 `form:"amount" binding:"required" label:"金额"`
	CouponId int64   `form:"coupon_id" label:"优惠卷id"`
	Code     string  `form:"code" binding:"required" label:"微信code"`
}

func (con RechargeController) Pay(c *gin.Context) {
	var data payRechargeControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	isAuth, code, msg, requestUser := Auth(GetToken(Token))

	if !isAuth {
		ay.Json{}.Msg(c, code, msg, gin.H{})
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

	amount := data.Amount - coupon.Amount

	//// 获取会员等级
	//isVip, vipLevel := CommonController{}.GetVip(requestUser)
	//if isVip {
	//	amount = data.Amount - coupon.Amount
	//} else {
	//
	//}
	v := strconv.FormatFloat(amount, 'g', -1, 64)

	isOrder, outTradeNo, order := CommonController{}.MakeOrder(1, 5, "充值"+v+"元", requestUser.Id, coupon.Id, amount, data.Amount, c.ClientIP(), "", 0, 0, 0)
	if isOrder == 0 {
		ay.Json{}.Msg(c, 400, outTradeNo, gin.H{})
		return
	}

	isC, msg, info := PayController{}.Wechat(data.Code, order.Des, order.OutTradeNo, c.ClientIP(), amount)

	if !isC {
		ay.Json{}.Msg(c, 400, msg, gin.H{})
		return
	}

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"info": info,
	})

}
