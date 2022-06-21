/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package models

import "dragonfly/ay"

type ConfigModel struct {
}

type Config struct {
	Id int64 `json:"id"`
	//RoomAmount         float64 `json:"room_amount"`
	UnsubscribeDesc    string  `json:"unsubscribe_desc"`
	VipPayDesc         string  `json:"vip_pay_desc"`
	RechargeDesc       string  `json:"recharge_desc"`
	VipPromotionDesc   string  `json:"vip_promotion_desc"`
	VipPayInstructions string  `json:"vip_pay_instructions"`
	KfLink             string  `json:"kf_link"`
	KfMobile           string  `json:"kf_mobile"`
	CouponType         string  `json:"coupon_type"`
	CouponName         string  `json:"coupon_name"`
	CouponSubName      string  `json:"coupon_sub_name"`
	CouponDes          string  `json:"coupon_des"`
	CouponAmount       float64 `json:"coupon_amount"`
	CouponEffective    int     `json:"coupon_effective"`
	GiveCard           int64   `json:"give_card"`
}

func (Config) TableName() string {
	return "d_config"
}

func (con ConfigModel) Get() (res Config) {
	ay.Db.First(&res, 1)
	return
}
