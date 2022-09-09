/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package models

import (
	"dragonfly/ay"
	"encoding/json"
	"strconv"
)

type UserRoomSubscribeModel struct {
}

type UserRoomSubscribe struct {
	BaseModel
	RoomId     int64  `json:"room_id"`
	Uid        int64  `json:"uid"`
	Ymd        int    `json:"ymd"`
	HourId     int    `json:"hour_id"`
	Status     int    `json:"status"`
	OutTradeNo string `json:"out_trade_no"`
}

func (UserRoomSubscribe) TableName() string {
	return "d_user_room_subscribe"
}

func (con UserRoomSubscribeModel) Add(outTradeNo string, hourId string, uid int64, roomId int64, ymd int) (bool, string) {
	type cv struct {
		Num  int      `json:"num"`
		Time []string `json:"time"`
		Ymd  int      `json:"ymd"`
	}
	var cc cv
	json.Unmarshal([]byte(hourId), &cc)
	for _, v := range cc.Time {
		hid, _ := strconv.Atoi(v)
		ay.Db.Create(&UserRoomSubscribe{
			RoomId:     roomId,
			Uid:        uid,
			Ymd:        ymd,
			HourId:     hid,
			OutTradeNo: outTradeNo,
		})
	}
	return true, "success"
}

func (con UserRoomSubscribeModel) Del(hourId string, uid int64, roomId int64) (bool, string) {
	type cv struct {
		Num  int      `json:"num"`
		Time []string `json:"time"`
		Ymd  int      `json:"ymd"`
	}
	var cc cv
	json.Unmarshal([]byte(hourId), &cc)
	for _, v := range cc.Time {
		var res UserRoomSubscribe
		ay.Db.Where("room_id = ? and uid = ? and ymd = ? and hour_id = ?", roomId, uid, cc.Ymd, v).First(&res)
		ay.Db.Delete(res)
	}
	return true, "success"
}
