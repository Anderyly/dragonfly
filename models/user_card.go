/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package models

import (
	"dragonfly/ay"
	"log"
	"time"
)

type UserCardModel struct {
}

type UserCard struct {
	BaseModel
	CardId      int64  `json:"card_id"`
	Uid         int64  `json:"uid"`
	UseHour     int    `json:"use_hour"`
	AllHour     int    `json:"all_hour"`
	EffectiveAt MyTime `json:"effective_at"`
	Status      int    `json:"status"`
}

func (UserCard) TableName() string {
	return "d_user_card"
}

func (con UserCardModel) Add(cardId int64, allHour int, uid int64, effectiveT int, amount float64) bool {
	effective := time.Now().Unix() + int64(effectiveT*3600*24)
	log.Println(effectiveT)
	log.Println(effective)
	stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Unix(effective, 0).Format("2006-01-02 15:04:05"), time.Local)
	if err := ay.Db.Create(&UserCard{
		CardId:      cardId,
		Uid:         uid,
		UseHour:     0,
		EffectiveAt: MyTime{Time: stamp},
		Status:      0,
		AllHour:     allHour,
	}).Error; err != nil {
		return false
	} else {
		var user User
		ay.Db.First(&user, uid)
		user.VipNum += amount
		ay.Db.Save(&user)
		return true
	}
}
