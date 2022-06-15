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

type SpellingAskController struct {
}

// List 列表
func (con SpellingAskController) List(c *gin.Context) {
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

	var list []models.ClassSpellingAsk

	var count int64
	res := ay.Db.Model(&models.ClassSpellingAsk{})

	res.Where("spelling_id = ?", data.Id)

	if data.Key != "" {
		res.Where("user_content like ? OR education_content like ?", "%"+data.Key+"%", "%"+data.Key+"%")
	}

	row := res

	row.Count(&count)

	res.Order("created_at desc").
		Limit(data.PageSize).
		Offset((data.Page - 1) * data.PageSize).
		Find(&list)

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"list":  list,
		"total": count,
	})
}

//type orderDetailForm struct {
//	Id int `form:"id"`
//}

// Detail 用户详情
func (con SpellingAskController) Detail(c *gin.Context) {
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

	var res models.ClassSpellingAsk

	ay.Db.First(&res, data.Id)

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"info": res,
	})
}

// Option 添加 编辑
func (con SpellingAskController) Option(c *gin.Context) {
	type optionForm struct {
		Id               int    `form:"id"`
		SpellingId       int64  `form:"spelling_id"`
		UserName         string `form:"user_name"`
		UserAvatar       string `form:"user_avatar"`
		UserContent      string `form:"user_content"`
		EducationName    string `form:"education_name"`
		EducationContent string `form:"education_content"`
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

	var res models.ClassSpellingAsk
	ay.Db.First(&res, data.Id)

	if data.Id != 0 {

		res.UserAvatar = data.UserAvatar
		res.UserName = data.UserName
		res.UserContent = data.UserContent
		res.EducationName = data.EducationName
		res.EducationContent = data.EducationContent

		if err := ay.Db.Save(&res).Error; err != nil {
			ay.Json{}.Msg(c, 400, "修改失败", gin.H{})
		} else {
			ay.Json{}.Msg(c, 200, "修改成功", gin.H{})
		}
	} else {
		ay.Db.Create(&models.ClassSpellingAsk{
			SpellingId:       data.SpellingId,
			UserName:         data.UserName,
			UserAvatar:       data.UserAvatar,
			UserContent:      data.UserContent,
			EducationName:    data.EducationName,
			EducationContent: data.EducationContent,
		})
		ay.Json{}.Msg(c, 200, "创建成功", gin.H{})

	}

}

func (con SpellingAskController) Delete(c *gin.Context) {
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
		var res models.ClassSpellingAsk
		ay.Db.Delete(&res, v)
	}

	ay.Json{}.Msg(c, 200, "删除成功", gin.H{})
}

func (con SpellingAskController) Upload(c *gin.Context) {

	code, msg, name := Upload(c, "spell_ask")

	if code != 200 {
		ay.Json{}.Msg(c, 400, msg, gin.H{})
	} else {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"url":  msg,
			"name": name,
		})
	}
}
