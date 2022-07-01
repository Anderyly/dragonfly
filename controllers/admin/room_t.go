/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package admin

import (
	"dragonfly/ay"
	"dragonfly/controllers"
	"dragonfly/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"strings"
)

type RoomTController struct {
}

// List 列表
func (con RoomTController) List(c *gin.Context) {
	type listForm struct {
		T      int    `form:"t"`
		RoomId string `form:"room_id"`
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

	var list []models.RoomSubscribe
	res := ay.Db.Model(&models.RoomSubscribe{})

	res.Where("room_id = ? AND ymd = ?", data.RoomId, data.T)

	res.Order("hour_id asc").
		Select("id,hour_id,amount").
		Find(&list)

	var room models.Room
	ay.Db.First(&room, data.RoomId).Select("amount")

	var h []gin.H

	for i := 1; i <= 24; i++ {
		j := 0
		for _, v := range list {
			if i == v.HourId {
				h = append(h, gin.H{
					"id":      v.Id,
					"amount":  strconv.FormatFloat(v.Amount, 'f', -1, 64),
					"hour_id": v.HourId,
				})
				j = 1
			}
		}
		if j == 0 {
			h = append(h, gin.H{
				"id":      0,
				"amount":  strconv.FormatFloat(room.Amount, 'f', -1, 64),
				"hour_id": i,
			})
		}
	}

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"list": h,
	})
}

//type orderDetailForm struct {
//	Id int `form:"id"`
//}

// Detail 用户详情
func (con RoomTController) Detail(c *gin.Context) {
	type detailForm struct {
		Id string `form:"id"`
	}
	var data detailForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	if Auth() == false {
		ay.Json{}.Msg(c, 401, "请登入", gin.H{})
		return
	}

	var res models.Room

	ay.Db.First(&res, data.Id)

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"info": res,
	})
}

// Option 添加 编辑
func (con RoomTController) Option(c *gin.Context) {
	type optionForm struct {
		Id     int64  `form:"id"`
		RoomId int64  `form:"room_id"`
		T      string `form:"t"`
		Ymd    int    `form:"ymd"`
	}

	type tt struct {
		Amount string `json:"amount"`
		HourId int    `json:"hour_id"`
		Id     int    `json:"id"`
	}

	var data optionForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	if Auth() == false {
		ay.Json{}.Msg(c, 401, "请登入", gin.H{})
		return
	}

	var room models.Room
	ay.Db.First(&room, data.RoomId)

	var ttt []tt
	json.Unmarshal([]byte(data.T), &ttt)
	log.Println(ttt)
	for _, v := range ttt {
		amount, _ := strconv.ParseFloat(v.Amount, 64)
		if amount == room.Amount {
			continue
		}
		var res models.RoomSubscribe
		ay.Db.Where("room_id = ? and ymd = ? and hour_id = ?", data.RoomId, data.Ymd, v.HourId).First(&res)
		if res.Id != 0 {
			log.Println(v.Amount)
			res.Amount = amount
			ay.Db.Save(&res)
		} else {
			ay.Db.Create(&models.RoomSubscribe{
				RoomId: data.RoomId,
				Ymd:    data.Ymd,
				HourId: v.HourId,
				Amount: amount,
			})
		}
	}
	ay.Json{}.Msg(c, 200, "修改成功", gin.H{})

}

func (con RoomTController) Delete(c *gin.Context) {
	type deleteForm struct {
		Id string `form:"id"`
	}
	var data deleteForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	if Auth() == false {
		ay.Json{}.Msg(c, 401, "请登入", gin.H{})
		return
	}

	idArr := strings.Split(data.Id, ",")

	for _, v := range idArr {
		var res models.Room
		ay.Db.First(&res, v)
		controllers.ControlServer{}.DelDept(res.DeptId)
		ay.Db.Delete(&res, v)
	}

	ay.Json{}.Msg(c, 200, "删除成功", gin.H{})
}

func (con RoomTController) Upload(c *gin.Context) {

	code, msg, name := Upload(c, "room")

	if code != 200 {
		ay.Json{}.Msg(c, 400, msg, gin.H{})
	} else {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"url":  msg,
			"name": name,
		})
	}
}
