/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package admin

import (
	"dragonfly/ay"
	"dragonfly/models"
	"github.com/gin-gonic/gin"
)

type BasicController struct {
}

// Detail 用户详情
func (con BasicController) Detail(c *gin.Context) {
	type detailForm struct {
		Id string `form:"id"`
	}
	var data detailForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	if Auth() == false {
		ay.Json{}.Msg(c, 401, "请登入", gin.H{})
		return
	}

	var res models.Config

	ay.Db.First(&res, 1)

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"info": res,
	})
}

// Option 添加 编辑
func (con BasicController) Option(c *gin.Context) {
	type optionForm struct {
		Id                 int     `form:"id"`
		RoomAmount         float64 `form:"room_amount"`
		UnsubscribeDesc    string  `form:"unsubscribe_desc"`
		VipPayDesc         string  `form:"vip_pay_desc"`
		RechargeDesc       string  `form:"recharge_desc"`
		VipPromotionDesc   string  `form:"vip_promotion_desc"`
		VipPayInstructions string  `form:"vip_pay_instructions"`
		KfLink             string  `form:"kf_link"`
		KfMobile           string  `form:"kf_mobile"`

		CouponType      string  `form:"coupon_type"`
		CouponName      string  `form:"coupon_name"`
		CouponSubName   string  `form:"coupon_sub_name"`
		CouponDes       string  `form:"coupon_des"`
		CouponAmount    float64 `form:"coupon_amount"`
		CouponEffective int     `form:"coupon_effective"`
		GiveCard        int64   `form:"give_card"`
	}
	var data optionForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	if Auth() == false {
		ay.Json{}.Msg(c, 401, "请登入", gin.H{})
		return
	}

	var res models.Config
	ay.Db.First(&res, 1)

	res.RechargeDesc = data.RechargeDesc
	res.KfMobile = data.KfMobile
	res.KfLink = data.KfLink
	res.UnsubscribeDesc = data.UnsubscribeDesc
	res.VipPayDesc = data.VipPayDesc
	res.VipPromotionDesc = data.VipPromotionDesc
	res.VipPayInstructions = data.VipPayInstructions

	res.CouponType = data.CouponType
	res.CouponName = data.CouponName
	res.CouponSubName = data.CouponSubName
	res.CouponDes = data.CouponDes
	res.CouponAmount = data.CouponAmount
	res.CouponEffective = data.CouponEffective
	res.GiveCard = data.GiveCard

	if err := ay.Db.Save(&res).Error; err != nil {
		ay.Json{}.Msg(c, 400, "修改失败", gin.H{})
	} else {
		ay.Json{}.Msg(c, 200, "修改成功", gin.H{})
	}

}
