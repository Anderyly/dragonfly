/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package models

import "dragonfly/ay"

type UserModel struct {
}

type User struct {
	BaseModel
	Openid      string  `json:"openid"`
	Avatar      string  `gorm:"column:avatar" json:"avatar"`
	NickName    string  `gorm:"column:nickname" json:"nickname"`
	Amount      float64 `gorm:"column:amount" json:"amount"`
	EffectiveAt MyTime  `gorm:"column:effective_at" json:"effective_at"`
	VipNum      float64 `gorm:"column:vip_num" json:"vip_num"`
}

func (User) TableName() string {
	return "d_user"
}

func (con UserModel) GetOpenid(openid string) (res User) {
	ay.Db.Where("openid = ?", openid).First(&res)
	return
}
