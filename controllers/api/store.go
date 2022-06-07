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

type StoreController struct {
}

type listStoreControllerForm struct {
	Page int `form:"page" binding:"required" label:"页码"`
}

func (con StoreController) List(c *gin.Context) {
	var data listStoreControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	page := data.Page - 1

	var list []models.BranchStore
	ay.Db.Select("id,name,label,image,address").
		Offset(page * Limit).
		Limit(Limit).
		Find(&list)

	var r []gin.H

	for _, v := range list {
		imageArr := strings.Split(v.Image, ",")

		var count int64
		ay.Db.Model(models.Room{}).Where("store_id = ?", v.Id).Count(&count)

		r = append(r, gin.H{
			"id":       v.Id,
			"name":     v.Name,
			"label":    v.Label,
			"image":    ay.Yaml.GetString("domain") + imageArr[0],
			"address":  v.Address,
			"room_num": count,
		})
	}

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"list": r,
	})

}

type detailStoreControllerForm struct {
	Id int `form:"id" binding:"required" label:"分店id"`
}

func (con StoreController) Detail(c *gin.Context) {
	var data detailStoreControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	var res models.BranchStore
	ay.Db.Where("id = ?", data.Id).First(&res)

	if res.Id == 0 {
		ay.Json{}.Msg(c, 400, "分店不存在", gin.H{})
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
		"address":   res.Address,
		"mobile":    res.Mobile,
		"image":     image,
		"introduce": res.Introduce,
		"attention": res.Attention,
		"video":     video,
	})

}
