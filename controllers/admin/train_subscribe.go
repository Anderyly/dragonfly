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
)

type TrainSubscribeController struct {
}

// List 列表
func (con TrainSubscribeController) List(c *gin.Context) {
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
		CreatedAt models.MyTime `json:"created_at"`
	}
	var list []r

	var count int64
	res := ay.Db.Table("d_train_log").
		Select("d_user.avatar,d_user.nickname,d_train_log.created_at").
		Joins("join d_user on d_train_log.uid=d_user.id")

	res.Where("d_train_log.train_id = ?", data.Id)

	row := res

	row.Count(&count)

	res.Order("d_train_log.created_at desc").
		Limit(data.PageSize).
		Offset((data.Page - 1) * data.PageSize).
		Find(&list)

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"list":  list,
		"total": count,
	})
}
