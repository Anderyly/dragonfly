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

type LevelController struct {
}

// List 列表
func (con LevelController) List(c *gin.Context) {

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

	var list []models.VipLevel

	var count int64
	res := ay.Db
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
func (con LevelController) Detail(c *gin.Context) {
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

	var user models.VipLevel

	ay.Db.First(&user, data.Id)

	ay.Json{}.Msg(c, 200,
		"success", gin.H{
			"info": user,
		})
}

// Option 添加 编辑
func (con LevelController) Option(c *gin.Context) {

	type optionForm struct {
		Id            int64   `form:"id"`
		Name          string  `form:"name"`
		Num           int     `form:"num"`
		Subscribe     float64 `form:"subscribe"`
		SpellClass    float64 `form:"spell_class"`
		Chanel        int     `form:"chanel"`
		AdvanceChanel int     `form:"advance_chanel"`
		GiveCard      string  `form:"give_card"`
		GiveCardId    int64   `form:"give_card_id"`
		GiveCoupon    int     `form:"give_coupon"`
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

	var res models.VipLevel
	ay.Db.First(&res, data.Id)
	if res.Id != 0 {
		res.Subscribe = data.Subscribe
		res.Name = data.Name
		res.Num = data.Num
		res.AdvanceChanel = data.AdvanceChanel
		res.SpellClass = data.SpellClass
		res.Chanel = data.Chanel
		res.GiveCoupon = data.GiveCoupon
		res.GiveCard = data.GiveCard
		res.GiveCardId = data.GiveCardId

		ay.Db.Save(&res)
		ay.Json{}.Msg(c, 200, "修改成功", gin.H{})
	} else {
		ay.Db.Create(&models.VipLevel{
			Name:          data.Name,
			Num:           data.Num,
			Subscribe:     data.Subscribe,
			SpellClass:    data.SpellClass,
			Chanel:        data.Chanel,
			AdvanceChanel: data.AdvanceChanel,
			GiveCard:      data.GiveCard,
			GiveCoupon:    data.GiveCoupon,
			GiveCardId:    data.GiveCardId,
		})
		ay.Json{}.Msg(c, 200, "创建成功", gin.H{})
	}

}

// Delete 删除
func (con LevelController) Delete(c *gin.Context) {
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
		var user models.VipLevel
		ay.Db.Delete(&user, v)
	}

	ay.Json{}.Msg(c, 200,
		"删除成功", gin.H{})
}
