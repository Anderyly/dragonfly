/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package models

type VipModel struct {
}

type Vip struct {
	BaseModel
	Name      string  `json:"name"`
	Amount    float64 `json:"amount"`
	OldAmount float64 `json:"old_amount"`
	Discount  float64 `json:"discount"`
	Sort      int     `json:"sort"`
}

func (Vip) TableName() string {
	return "d_vip"
}
