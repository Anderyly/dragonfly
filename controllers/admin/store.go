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

type StoreController struct {
}

// List 列表
func (con StoreController) List(c *gin.Context) {
	type listForm struct {
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
		Key      string `form:"key"`
		Status   string `form:"status"`
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

	var list []models.BranchStore

	var count int64
	res := ay.Db.Model(&models.BranchStore{})

	if data.Key != "" {
		res.Where("name like ?", "%"+data.Key+"%")
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
func (con StoreController) Detail(c *gin.Context) {
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

	var res models.BranchStore

	ay.Db.First(&res, data.Id)

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"info": res,
	})
}

// Option 添加 编辑
func (con StoreController) Option(c *gin.Context) {
	type optionForm struct {
		Id        int    `form:"id"`
		Name      string `form:"name"`
		Attention string `form:"attention"`
		Label     string `form:"label"`
		Introduce string `form:"introduce"`
		Image     string `form:"image"`
		Video     string `form:"video"`
		Mobile    string `form:"mobile"`
		Address   string `form:"address"`
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

	var res models.BranchStore
	ay.Db.First(&res, data.Id)

	if data.Id != 0 {

		res.Image = strings.ReplaceAll(data.Image, ay.Yaml.GetString("domain"), "")
		res.Name = data.Name
		res.Introduce = data.Introduce
		res.Attention = data.Attention
		res.Label = data.Label
		res.Video = data.Video
		res.Mobile = data.Mobile
		res.Address = data.Address

		if err := ay.Db.Save(&res).Error; err != nil {
			ay.Json{}.Msg(c, 400, "修改失败", gin.H{})
		} else {
			ay.Json{}.Msg(c, 200, "修改成功", gin.H{})
		}
	} else {
		ay.Db.Create(&models.BranchStore{
			Name:      data.Name,
			Address:   data.Address,
			Mobile:    data.Mobile,
			Introduce: data.Introduce,
			Attention: data.Attention,
			Image:     strings.ReplaceAll(data.Image, ay.Yaml.GetString("domain"), ""),
			Video:     data.Video,
			Label:     data.Label,
		})
		ay.Json{}.Msg(c, 200, "创建成功", gin.H{})

	}

}

func (con StoreController) Delete(c *gin.Context) {
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
		var res models.BranchStore
		ay.Db.Delete(&res, v)
	}

	ay.Json{}.Msg(c, 200, "删除成功", gin.H{})
}

func (con StoreController) Upload(c *gin.Context) {

	code, msg, name := Upload(c, "store")

	if code != 200 {
		ay.Json{}.Msg(c, 400, msg, gin.H{})
	} else {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"url":  msg,
			"name": name,
		})
	}
}

// GetAll 获取所有分店
func (con StoreController) GetAll(c *gin.Context) {
	if Auth() == false {
		ay.Json{}.Msg(c, 401, "请登入", gin.H{})
		return
	}

	type returnList struct {
		Label string `gorm:"column:name" json:"label"`
		Value int64  `gorm:"column:id" json:"value"`
	}

	var list []returnList
	ay.Db.Model(models.BranchStore{}).Find(&list)

	ay.Json{}.Msg(c, 200,
		"success", gin.H{
			"list": list,
		})

}
