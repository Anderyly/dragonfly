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
	"strings"
	"time"
)

type AccountCardController struct {
}

// List 列表
func (con AccountCardController) List(c *gin.Context) {

	type listForm struct {
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
		Key      string `form:"key"`
		Uid      int64  `form:"uid"`
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
		Id          int64         `json:"id"`
		UseHour     float64       `json:"use_hour"`
		EffectiveAt models.MyTime `json:"effective_at"`
		Status      int           `json:"status"`
		CreatedAt   models.MyTime `json:"created_at"`
		Name        string        `json:"name"`
		Type        int           `json:"type"`
		AllHour     float64       `json:"all_hour"`
	}

	var list []r

	var count int64
	res := ay.Db.Table("d_user_card").
		Select("d_user_card.id,d_user_card.use_hour,d_user_card.effective_at,d_user_card.created_at,d_user_card.status,d_card.name,d_card.type,d_card.all_hour").
		Joins("join d_card on d_user_card.card_id=d_card.id")

	res.Where("d_user_card.uid = ?", data.Uid)

	row := res

	row.Count(&count)

	res.Order("d_user_card.created_at desc").
		Limit(data.PageSize).
		Offset((data.Page - 1) * data.PageSize).
		Find(&list)

	ay.Json{}.Msg(c, 200,
		"success", gin.H{
			"list":  list,
			"total": count,
		})
}

// Detail 详情
func (con AccountCardController) Detail(c *gin.Context) {
	type detailForm struct {
		Id int `form:"id"`
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

	var user models.UserCard

	ay.Db.First(&user, data.Id)

	ay.Json{}.Msg(c, 200,
		"success", gin.H{
			"info": user,
		})
}

// Option 添加 编辑
func (con AccountCardController) Option(c *gin.Context) {

	type optionForm struct {
		Id          int64  `form:"id"`
		UseHour     int    `form:"use_hour"`
		Uid         int64  `form:"uid"`
		CardId      int64  `form:"card_id"`
		Status      int    `form:"status"`
		EffectiveAt string `form:"effective_at"`
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

	var res models.UserCard
	ay.Db.First(&res, data.Id)
	stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", data.EffectiveAt, time.Local)

	if res.Id != 0 {

		res.EffectiveAt = models.MyTime{Time: stamp}

		res.Status = data.Status
		res.UseHour = data.UseHour
		res.CardId = data.CardId
		ay.Db.Save(&res)
		ay.Json{}.Msg(c, 200, "修改成功", gin.H{})
	} else {
		var card models.Card
		ay.Db.First(&card, data.CardId)
		ay.Db.Create(&models.UserCard{
			BaseModel:   models.BaseModel{},
			CardId:      data.CardId,
			Uid:         data.Uid,
			UseHour:     data.UseHour,
			AllHour:     card.AllHour,
			EffectiveAt: models.MyTime{Time: stamp},
			Status:      data.Status,
		})
		ay.Json{}.Msg(c, 200, "创建成功", gin.H{})
	}

}

// Delete 删除
func (con AccountCardController) Delete(c *gin.Context) {
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
		var user models.UserCard
		ay.Db.Delete(&user, v)
	}

	ay.Json{}.Msg(c, 200,
		"删除成功", gin.H{})
}
