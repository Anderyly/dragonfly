/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package models

type RoomSubscribeModel struct {
}

type RoomSubscribe struct {
	BaseModel
	RoomId int64   `json:"room_id"`
	Ymd    int     `json:"ymd"`
	HourId int     `json:"hour_id"`
	Amount float64 `json:"amount"`
}

func (RoomSubscribe) TableName() string {
	return "d_room_subscribe"
}
