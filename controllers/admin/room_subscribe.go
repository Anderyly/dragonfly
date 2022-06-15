/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package admin

import (
	"dragonfly/ay"
	"dragonfly/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

type RoomSubscribeController struct {
}

// List 列表
func (con RoomSubscribeController) List(c *gin.Context) {
	type listForm struct {
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
		Key      string `form:"key"`
		Status   string `form:"status"`
		Type     string `form:"type"`
		Id       string `form:"id"`
	}
	var data listForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	if Auth() == false {
		ay.Json{}.Msg(c, 401, "请登入", gin.H{})
		return
	}

	type r struct {
		Avatar    string        `json:"avatar"`
		Nickname  string        `json:"nickname"`
		Ymd       string        `json:"ymd"`
		HourId    string        `json:"hour_id"`
		CreatedAt models.MyTime `json:"created_at"`
	}
	var list []r

	var count int64
	res := ay.Db.Table("d_user_room_subscribe").
		Select("d_user.avatar,d_user.nickname,d_user_room_subscribe.ymd,d_user_room_subscribe.hour_id,d_user_room_subscribe.created_at").
		Joins("join d_user on d_user_room_subscribe.uid=d_user.id")

	res.Where("d_user_room_subscribe.room_id = ?", data.Id)

	row := res

	row.Count(&count)

	res.Order("d_user_room_subscribe.created_at desc").
		Limit(data.PageSize).
		Offset((data.Page - 1) * data.PageSize).
		Find(&list)

	sjd := [24]string{
		"00:00-01:00",
		"01:00-02:00",
		"02:00-03:00",
		"03:00-04:00",
		"04:00-05:00",
		"05:00-06:00",
		"06:00-07:00",
		"07:00-08:00",
		"08:00-09:00",
		"09:00-10:00",
		"10:00-11:00",
		"11:00-12:00",
		"12:00-13:00",
		"13:00-14:00",
		"14:00-15:00",
		"15:00-16:00",
		"16:00-17:00",
		"17:00-18:00",
		"18:00-19:00",
		"19:00-20:00",
		"20:00-21:00",
		"21:00-22:00",
		"22:00-23:00",
		"23:00-00:00",
	}

	for k, v := range list {

		i, _ := strconv.Atoi(v.HourId)
		list[k].HourId = sjd[i]
		list[k].Ymd = v.Ymd[0:4] + "-" + v.Ymd[4:6] + "-" + v.Ymd[6:8]
	}

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"list":  list,
		"total": count,
	})
}
