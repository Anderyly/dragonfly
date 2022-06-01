/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package models

type VipLevelModel struct {
}

type VipLevel struct {
	BaseModel
	Name          string  `json:"name"`
	Num           int     `json:"num"`
	Subscribe     float64 `json:"subscribe"`
	SpellClass    float64 `json:"spell_class"`
	Chanel        int     `json:"chanel"`
	AdvanceChanel int     `json:"advance_chanel"`
	GiveCard      int     `json:"give_card"`
	GiveCoupon    int     `json:"give_coupon"`
}

func (VipLevel) TableName() string {
	return "d_vip_level"
}
