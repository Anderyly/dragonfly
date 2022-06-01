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
)

type EducationController struct {
}

type listEducationControllerForm struct {
	StoreId int `form:"store_id" binding:"required" label:"分店id"`
	Page    int `form:"page" binding:"required" label:"页码"`
}

func (con EducationController) List(c *gin.Context) {
	var data listEducationControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	page := data.Page - 1

	var list []models.Education
	ay.Db.Select("id,name,avatar,label,wechat,introduce").
		Where("store_id = ?", data.StoreId).
		Offset(page * Limit).
		Limit(Limit).Debug().
		Find(&list)

	var r []gin.H

	for _, v := range list {

		r = append(r, gin.H{
			"id":        v.Id,
			"name":      v.Name,
			"wechat":    v.Wechat,
			"label":     v.Label,
			"introduce": v.Introduce,
			"avatar":    ay.Yaml.GetString("domain") + v.Avatar,
		})
	}

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"list": r,
	})

}

type detailEducationControllerForm struct {
	Id int `form:"id" binding:"required" label:"私教id"`
}

func (con EducationController) Detail(c *gin.Context) {
	var data detailEducationControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	var res models.Education
	ay.Db.Where("id = ?", data.Id).First(&res)

	if res.Id == 0 {
		ay.Json{}.Msg(c, 400, "私教不存在", gin.H{})
		return
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

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"id":        res.Id,
		"name":      res.Name,
		"sign":      res.Sign,
		"label":     res.Label,
		"image":     image,
		"introduce": res.Introduce,
		"amount":    res.Amount,
		"wechat":    res.Wechat,
		"avatar":    ay.Yaml.GetString("domain") + res.Avatar,
		"video":     video,
	})

}
