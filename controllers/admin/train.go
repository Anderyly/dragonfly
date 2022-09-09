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

type TrainController struct {
}

// List 列表
func (con TrainController) List(c *gin.Context) {
	type listForm struct {
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
		Key      string `form:"key"`
		Status   string `form:"status"`
		Type     string `form:"type"`
		StoreId  string `form:"store_id"`
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

	var list []models.Train

	var count int64
	res := ay.Db.Model(&models.Train{})

	if data.Key != "" {
		res.Where("name like ?", "%"+data.Key+"%")
	}
	if data.StoreId != "" {
		res.Where("store_id = ?", data.StoreId)
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
func (con TrainController) Detail(c *gin.Context) {
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

	var res models.Train

	ay.Db.First(&res, data.Id)

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"info": res,
	})
}

// Option 添加 编辑
func (con TrainController) Option(c *gin.Context) {
	type optionForm struct {
		Id          int     `form:"id"`
		StoreId     int64   `form:"store_id"`
		Room        string  `form:"room"`
		EducationId int64   `form:"education_id"`
		Name        string  `form:"name"`
		Address     string  `form:"address"`
		Mobile      string  `form:"mobile"`
		Introduce   string  `form:"introduce"`
		Attention   string  `form:"attention"`
		Arrange     string  `form:"arrange"`
		BackImage   string  `form:"back_image"`
		Image       string  `form:"image"`
		Video       string  `form:"video"`
		MinUser     int     `form:"min_user"`
		MaxUser     int     `form:"max_user"`
		StartDate   string  `form:"start_date"`
		Amount      float64 `form:"amount"`
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

	var res models.Train
	ay.Db.First(&res, data.Id)

	stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", data.StartDate, time.Local)

	if data.Id != 0 {

		res.StoreId = data.StoreId
		res.Image = strings.ReplaceAll(data.Image, ay.Yaml.GetString("domain"), "")
		res.Name = data.Name
		res.Introduce = data.Introduce
		res.Video = data.Video
		res.BackImage = data.BackImage
		res.Attention = data.Attention
		res.Mobile = data.Mobile
		res.MinUser = data.MinUser
		res.MaxUser = data.MaxUser
		res.Room = data.Room
		res.EducationId = data.EducationId
		res.Address = data.Address
		res.Arrange = data.Arrange
		res.StartDate = models.MyTime{Time: stamp}
		res.Amount = data.Amount

		if err := ay.Db.Save(&res).Error; err != nil {
			ay.Json{}.Msg(c, 400, "修改失败", gin.H{})
		} else {
			ay.Json{}.Msg(c, 200, "修改成功", gin.H{})
		}
	} else {
		ay.Db.Create(&models.Train{
			StoreId:     data.StoreId,
			Room:        data.Room,
			EducationId: data.EducationId,
			Name:        data.Name,
			Address:     data.Address,
			Mobile:      data.Mobile,
			Introduce:   data.Introduce,
			Attention:   data.Introduce,
			Arrange:     data.Arrange,
			BackImage:   data.BackImage,
			Image:       strings.ReplaceAll(data.Image, ay.Yaml.GetString("domain"), ""),
			Video:       data.Video,
			MinUser:     data.MinUser,
			MaxUser:     data.MaxUser,
			StartDate:   models.MyTime{Time: stamp},
			Amount:      data.Amount,
		})
		ay.Json{}.Msg(c, 200, "创建成功", gin.H{})

	}

}

func (con TrainController) Delete(c *gin.Context) {
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
		var res models.Train
		ay.Db.Delete(&res, v)
	}

	ay.Json{}.Msg(c, 200, "删除成功", gin.H{})
}

func (con TrainController) Upload(c *gin.Context) {

	code, msg, name := Upload(c, "train")

	if code != 200 {
		ay.Json{}.Msg(c, 400, msg, gin.H{})
	} else {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"url":  msg,
			"name": name,
		})
	}
}
