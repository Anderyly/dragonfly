/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package models

type CouponModel struct {
}

type Coupon struct {
	BaseModel
	Type        string  `json:"type"`
	Uid         int64   `json:"uid"`
	Name        string  `json:"name"`
	SubName     string  `json:"sub_name"`
	Des         string  `json:"des"`
	Amount      float64 `json:"amount"`
	EffectiveAt MyTime  `json:"effective_at"`
	Status      int     `json:"status"`
}

func (Coupon) TableName() string {
	return "d_user_coupon"
}
