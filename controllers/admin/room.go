/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package admin

import (
	"dragonfly/ay"
	"dragonfly/controllers"
	"dragonfly/models"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type RoomController struct {
}

// List 列表
func (con RoomController) List(c *gin.Context) {
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

	var list []models.Room

	var count int64
	res := ay.Db.Model(&models.Room{})

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
func (con RoomController) Detail(c *gin.Context) {
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

	var res models.Room

	ay.Db.First(&res, data.Id)

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"info": res,
	})
}

// Option 添加 编辑
func (con RoomController) Option(c *gin.Context) {
	type optionForm struct {
		Id        int     `form:"id"`
		StoreId   int64   `form:"store_id"`
		Type      int     `form:"type"`
		Name      string  `form:"name"`
		Address   string  `form:"address"`
		Mobile    string  `form:"mobile"`
		Introduce string  `form:"introduce"`
		Attention string  `form:"attention"`
		Guide     string  `form:"guide"`
		Cancel    string  `form:"cancel"`
		Image     string  `form:"image"`
		Video     string  `form:"video"`
		Size      float64 `form:"size"`
		Electric  string  `form:"electric"`
		ControlSn string  `form:"control_sn"`
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

	var res models.Room
	ay.Db.First(&res, data.Id)

	image := strings.ReplaceAll(data.Image, ay.Yaml.GetString("domain"), "")

	id := controllers.ControlServer{}.GetId(data.ControlSn)
	if id == 0 {
		ay.Json{}.Msg(c, 400, "设备sn错误", gin.H{})
		return
	}
	if data.Id != 0 {

		res.StoreId = data.StoreId
		res.Image = image
		res.Name = data.Name
		res.Introduce = data.Introduce
		res.Video = data.Video
		res.Attention = data.Attention
		res.Mobile = data.Mobile
		res.Address = data.Address
		res.Size = data.Size
		res.Guide = data.Guide
		res.Cancel = data.Cancel
		res.Type = data.Type
		res.Electric = data.Electric
		res.ControlSn = data.ControlSn
		res.ControlId = strconv.Itoa(id)

		if err := ay.Db.Save(&res).Error; err != nil {
			ay.Json{}.Msg(c, 400, "修改失败", gin.H{})
		} else {
			controllers.ControlServer{}.EditDept(data.Name, res.DeptId)
			ay.Json{}.Msg(c, 200, "修改成功", gin.H{})
		}
	} else {

		_, ids := controllers.ControlServer{}.AddDept(data.Name, "0")
		ay.Db.Create(&models.Room{
			StoreId:   data.StoreId,
			Type:      data.Type,
			Name:      data.Name,
			Address:   data.Address,
			Mobile:    data.Mobile,
			Introduce: data.Introduce,
			Attention: data.Attention,
			Guide:     data.Guide,
			Cancel:    data.Cancel,
			Image:     image,
			Video:     data.Video,
			Size:      data.Size,
			Electric:  data.Electric,
			ControlSn: data.ControlSn,
			ControlId: strconv.Itoa(id),
			DeptId:    strconv.Itoa(ids),
		})
		ay.Json{}.Msg(c, 200, "创建成功", gin.H{})

	}

}

func (con RoomController) Delete(c *gin.Context) {
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
		var res models.Room
		ay.Db.First(&res, v)
		controllers.ControlServer{}.DelDept(res.DeptId)
		ay.Db.Delete(&res, v)
	}

	ay.Json{}.Msg(c, 200, "删除成功", gin.H{})
}

func (con RoomController) Upload(c *gin.Context) {

	code, msg, name := Upload(c, "room")

	if code != 200 {
		ay.Json{}.Msg(c, 400, msg, gin.H{})
	} else {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"url":  msg,
			"name": name,
		})
	}
}
