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

type RoomController struct {
}

type listRoomControllerForm struct {
	StoreId int `form:"store_id" binding:"required" label:"分店id"`
	Page    int `form:"page" binding:"required" label:"页码"`
}

func (con RoomController) List(c *gin.Context) {
	var data listRoomControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	page := data.Page - 1

	var list []models.Room
	ay.Db.Select("id,name,size,image,address").
		Where("store_id = ?", data.StoreId).
		Offset(page * Limit).
		Limit(Limit).Debug().
		Find(&list)

	var r []gin.H

	for _, v := range list {
		imageArr := strings.Split(v.Image, ",")

		r = append(r, gin.H{
			"id":      v.Id,
			"name":    v.Name,
			"size":    v.Size,
			"label":   "火爆",
			"image":   ay.Yaml.GetString("domain") + imageArr[0],
			"address": v.Address,
		})
	}

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"list": r,
	})

}

type detailRoomControllerForm struct {
	Id int `form:"id" binding:"required" label:"分店id"`
}

func (con RoomController) Detail(c *gin.Context) {
	var data detailRoomControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	var res models.Room
	ay.Db.Where("id = ?", data.Id).First(&res)

	if res.Id == 0 {
		ay.Json{}.Msg(c, 400, "舞蹈室不存在", gin.H{})
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
		"guide":     res.Guide,
		"cancel":    res.Cancel,
		"video":     video,
	})

}
