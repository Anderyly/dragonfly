/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package models

type UserCardModel struct {
}

type UserCard struct {
	BaseModel
	CardId      int64  `json:"card_id"`
	Uid         int64  `json:"uid"`
	UseHour     int    `json:"use_hour"`
	EffectiveAt MyTime `json:"effective_at"`
	Status      int    `json:"status"`
}

func (UserCard) TableName() string {
	return "d_user_card"
}
