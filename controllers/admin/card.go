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
)

type CardController struct {
}

// List 列表
func (con CardController) List(c *gin.Context) {

	type listForm struct {
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
		Key      string `form:"key"`
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

	var list []models.Card

	var count int64
	res := ay.Db.Model(models.Card{})
	if data.Key != "" {
		res.Where("name like ?", "%"+data.Key+"%")
	}

	row := res

	row.Count(&count)

	res.Order("created_at desc").
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
func (con CardController) Detail(c *gin.Context) {
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

	var user models.Card

	ay.Db.First(&user, data.Id)

	ay.Json{}.Msg(c, 200,
		"success", gin.H{
			"info": user,
		})
}

// Option 添加 编辑
func (con CardController) Option(c *gin.Context) {

	type optionForm struct {
		Id            int64   `form:"id"`
		Name          string  `form:"name"`
		Type          int     `form:"type"`
		Amount        float64 `form:"amount"`
		AllHour       int     `form:"all_hour"`
		EffectiveText string  `form:"effective_text"`
		Effective     int     `form:"effective"`
		HourAmount    float64 `form:"hour_amount"`
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

	var res models.Card
	ay.Db.First(&res, data.Id)
	if res.Id != 0 {
		res.Amount = data.Amount
		res.Effective = data.Effective
		res.EffectiveText = data.EffectiveText
		res.AllHour = data.AllHour
		res.HourAmount = data.HourAmount
		res.Type = data.Type
		res.Name = data.Name

		ay.Db.Save(&res)
		ay.Json{}.Msg(c, 200, "修改成功", gin.H{})
	} else {
		ay.Db.Create(&models.Card{
			Type:          data.Type,
			Name:          data.Name,
			Amount:        data.Amount,
			HourAmount:    data.HourAmount,
			AllHour:       data.AllHour,
			EffectiveText: data.EffectiveText,
			Effective:     data.Effective,
		})
		ay.Json{}.Msg(c, 200, "创建成功", gin.H{})
	}

}

// Delete 删除
func (con CardController) Delete(c *gin.Context) {
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
		var user models.Card
		ay.Db.Delete(&user, v)
	}

	ay.Json{}.Msg(c, 200,
		"删除成功", gin.H{})
}

// GetAll 获取所有次卡
func (con CardController) GetAll(c *gin.Context) {
	if Auth() == false {
		ay.Json{}.Msg(c, 401, "请登入", gin.H{})
		return
	}

	type returnList struct {
		Label string `gorm:"column:name" json:"label"`
		Value int64  `gorm:"column:id" json:"value"`
	}

	var list []returnList
	ay.Db.Model(models.Card{}).Find(&list)

	ay.Json{}.Msg(c, 200,
		"success", gin.H{
			"list": list,
		})

}
