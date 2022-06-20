/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package models

import (
	"dragonfly/ay"
	"time"
)

type UserModel struct {
}

type User struct {
	BaseModel
	Openid        string  `json:"openid"`
	Avatar        string  `gorm:"column:avatar" json:"avatar"`
	NickName      string  `gorm:"column:nickname" json:"nickname"`
	Amount        float64 `gorm:"column:amount" json:"amount"`
	EffectiveAt   MyTime  `gorm:"column:effective_at" json:"effective_at"`
	VipNum        float64 `gorm:"column:vip_num" json:"vip_num"`
	ControlUserId string  `json:"control_user_id"`
}

func (User) TableName() string {
	return "d_user"
}

func (con UserModel) GetOpenid(openid string) (res User) {
	ay.Db.Where("openid = ?", openid).First(&res)
	return
}

func (con UserModel) SetVip(user User, t int, amount float64) bool {
	effective := time.Now().Unix()
	if user.EffectiveAt.Unix() >= time.Now().Unix() {
		effective = user.EffectiveAt.Unix()
	}
	effective += int64(t * 3600 * 24)
	stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Unix(effective, 0).Format("2006-01-02 15:04:05"), time.Local)
	user.EffectiveAt = MyTime{Time: stamp}
	user.VipNum += amount
	if err := ay.Db.Save(&user).Error; err != nil {
		return false
	} else {
		return true
	}
}
