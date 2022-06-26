/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package api

import (
	"dragonfly/ay"
	"dragonfly/models"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type SpellingController struct {
}

type listSpellingControllerForm struct {
	StoreId int `form:"store_id" binding:"required" label:"分店id"`
	Page    int `form:"page" binding:"required" label:"页码"`
}

func (con SpellingController) List(c *gin.Context) {
	var data listSpellingControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	page := data.Page - 1

	var list []models.ClassSpelling
	ay.Db.Select("id,name,start_date,back_image").
		Where("store_id = ?", data.StoreId).
		Offset(page * Limit).
		Limit(Limit).
		Find(&list)

	var r []gin.H

	for _, v := range list {

		r = append(r, gin.H{
			"id":         v.Id,
			"name":       v.Name,
			"back_image": ay.Yaml.GetString("domain") + v.BackImage,
			"start_date": v.StartDate.Unix(),
		})
	}

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"list": r,
	})

}

type detailSpellingControllerForm struct {
	Id int `form:"id" binding:"required" label:"拼课id"`
}

func (con SpellingController) Detail(c *gin.Context) {
	var data detailSpellingControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	var res models.ClassSpelling
	ay.Db.Where("id = ?", data.Id).First(&res)

	if res.Id == 0 {
		ay.Json{}.Msg(c, 400, "拼课不存在", gin.H{})
		return
	}

	userSubscribe := 0
	status := 0

	_, code, _, requestUser := Auth(GetToken(Token))

	if code == 200 {
		var count int64
		ay.Db.Model(&models.ClassSpellingLog{}).Where("spelling_id = ? AND uid = ?", data.Id, requestUser.Id).Count(&count)
		userSubscribe = int(count)
		if count > 0 {
			userSubscribe = 1
		}
	}

	// 预约人数
	var count int64
	ay.Db.Model(&models.ClassSpellingLog{}).Where("spelling_id = ?", data.Id).Count(&count)
	if int(count) >= res.MaxUser {
		status = 2
	} else {
		status = 1
	}

	if time.Now().Unix() > res.StartDate.Unix() {
		status = 3
	}

	imageArr := strings.Split(res.Image, ",")

	var image []string

	for _, v := range imageArr {
		image = append(image, ay.Yaml.GetString("domain")+v)
	}

	video := ""
	if res.Video != "" {
		video = ay.Yaml.GetString("domain") + res.Video
	}

	var education models.Education
	ay.Db.First(&education, res.EducationId)

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"id":                  res.Id,
		"education_avatar":    ay.Yaml.GetString("domain") + education.Avatar,
		"education_name":      education.Name,
		"education_label":     education.Label,
		"education_introduce": education.Introduce,
		"education_wechat":    education.Wechat,
		"name":                res.Name,
		"address":             res.Address,
		"mobile":              res.Mobile,
		"image":               image,
		"introduce":           res.Introduce,
		"attention":           res.Attention,
		"video":               video,
		"min_user":            res.MinUser,
		"max_user":            res.MaxUser,
		"room":                res.Room,
		"start_date":          res.StartDate.Unix(),
		"arrange":             res.Arrange,
		"user_subscribe":      userSubscribe,
		"status":              status,
		"count":               res.MaxUser - int(count),
	})

}

type subscribeSpellingControllerForm struct {
	Id int `form:"id" binding:"required" label:"拼课id"`
}

// Subscribe 预约、取消预约
func (con SpellingController) Subscribe(c *gin.Context) {
	var data subscribeSpellingControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	isLogin, code, msg, requestUser := Auth(GetToken(Token))

	if !isLogin {
		ay.Json{}.Msg(c, code, msg, gin.H{})
		return
	}

	var res models.ClassSpelling
	ay.Db.Where("id = ?", data.Id).First(&res)

	if res.Id == 0 {
		ay.Json{}.Msg(c, 400, "拼课不存在", gin.H{})
		return
	}

	var spellingLog models.ClassSpellingLog
	ay.Db.Where("spelling_id = ? AND uid = ?", data.Id, requestUser.Id).First(&spellingLog)

	if spellingLog.Id != 0 {
		if err := ay.Db.Delete(&spellingLog).Error; err != nil {
			ay.Json{}.Msg(c, 400, "取消预约失败", gin.H{})
		} else {
			ay.Json{}.Msg(c, 200, "已取消预约", gin.H{})
		}
	} else {
		// 预约人数
		var count int64
		ay.Db.Model(&models.ClassSpellingLog{}).Where("spelling_id = ?", data.Id).Count(&count)
		if int(count) >= res.MaxUser {
			ay.Json{}.Msg(c, 400, "预约人数已满", gin.H{})
			return
		}
		if err := ay.Db.Create(&models.ClassSpellingLog{
			SpellingId: res.Id,
			Uid:        requestUser.Id,
		}).Error; err != nil {
			ay.Json{}.Msg(c, 400, "预约失败", gin.H{})
		} else {
			ay.Json{}.Msg(c, 200, "已预约成功", gin.H{})
		}
	}
}

type askSpellingControllerForm struct {
	SpellingId int `form:"spelling_id" binding:"required" label:"拼课id"`
	Page       int `form:"page" binding:"required" label:"页码"`
}

func (con SpellingController) Ask(c *gin.Context) {
	var data askSpellingControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	page := data.Page - 1

	var list []models.ClassSpellingAsk
	ay.Db.Where("spelling_id = ?", data.SpellingId).
		Offset(page * Limit).
		Limit(Limit).
		Order("created_at desc").
		Find(&list)

	var r []gin.H

	for _, v := range list {

		r = append(r, gin.H{
			"user_name":         v.UserName,
			"user_avatar":       ay.Yaml.GetString("domain") + v.UserAvatar,
			"user_content":      v.UserContent,
			"education_name":    v.EducationName,
			"education_content": v.EducationContent,
			"created_at":        v.CreatedAt.Unix(),
		})
	}

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"list": r,
	})

}
