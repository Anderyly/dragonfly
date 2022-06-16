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

type OrderController struct {
}

// List 列表
func (con OrderController) List(c *gin.Context) {
	type listForm struct {
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
		Key      string `form:"key"`
		Status   string `form:"status"`
		Type     string `form:"type"`
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

	type r struct {
		Nickname string `json:"nickname"`
		models.Order
	}
	var list []r

	var count int64
	res := ay.Db.Table("d_order").
		Select("d_user.nickname,d_order.*").
		Joins("join d_user on d_order.uid=d_user.id")

	if data.Key != "" {
		res.Where("d_order.out_trade_no like ? or d_order.trade_no like ? or d_user.nickname like ?", "%"+data.Key+"%", "%"+data.Key+"%", "%"+data.Key+"%")
	}
	if data.Type != "" {
		res.Where("d_order.type = ?", data.Type)
	}
	if data.Status != "" {
		res.Where("d_order.status = ?", data.Status)
	}

	row := res

	row.Count(&count)

	res.Order("d_order.created_at desc").
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
func (con OrderController) Detail(c *gin.Context) {
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

	var res models.Education

	ay.Db.First(&res, data.Id)

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"info": res,
	})
}

// Option 添加 编辑
func (con OrderController) Option(c *gin.Context) {
	type optionForm struct {
		Id        int     `form:"id"`
		StoreId   int64   `form:"store_id"`
		Name      string  `form:"name"`
		Avatar    string  `form:"avatar"`
		Label     string  `form:"label"`
		Introduce string  `form:"introduce"`
		Amount    float64 `form:"amount"`
		Sign      string  `form:"sign"`
		Wechat    string  `form:"wechat"`
		Image     string  `form:"image"`
		Video     string  `form:"video"`
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

	var res models.Education
	ay.Db.First(&res, data.Id)

	if data.Id != 0 {

		res.StoreId = data.StoreId
		res.Image = strings.ReplaceAll(data.Image, ay.Yaml.GetString("domain"), "")
		res.Amount = data.Amount
		res.Name = data.Name
		res.Avatar = data.Avatar
		res.Wechat = data.Wechat
		res.Sign = data.Sign
		res.Introduce = data.Introduce
		res.Label = data.Label
		res.Video = data.Video

		if err := ay.Db.Save(&res).Error; err != nil {
			ay.Json{}.Msg(c, 400, "修改失败", gin.H{})
		} else {
			ay.Json{}.Msg(c, 200, "修改成功", gin.H{})
		}
	} else {
		ay.Db.Create(&models.Education{
			StoreId:   data.StoreId,
			Name:      data.Name,
			Avatar:    data.Avatar,
			Label:     data.Label,
			Introduce: data.Introduce,
			Amount:    data.Amount,
			Sign:      data.Sign,
			Wechat:    data.Wechat,
			Image:     strings.ReplaceAll(data.Image, ay.Yaml.GetString("domain"), ""),
			Video:     data.Video,
		})
		ay.Json{}.Msg(c, 200, "创建成功", gin.H{})

	}

}

func (con OrderController) Delete(c *gin.Context) {
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
		var res models.Education
		ay.Db.Delete(&res, v)
	}

	ay.Json{}.Msg(c, 200, "删除成功", gin.H{})
}

func (con OrderController) Upload(c *gin.Context) {

	code, msg, name := Upload(c, "education")

	if code != 200 {
		ay.Json{}.Msg(c, 400, msg, gin.H{})
	} else {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"url":  msg,
			"name": name,
		})
	}
}

// GetAll 获取所有私教
func (con OrderController) GetAll(c *gin.Context) {
	if Auth() == false {
		ay.Json{}.Msg(c, 401, "请登入", gin.H{})
		return
	}

	type returnList struct {
		Label string `gorm:"column:name" json:"label"`
		Value int64  `gorm:"column:id" json:"value"`
	}

	var list []returnList
	ay.Db.Model(models.Education{}).Find(&list)

	ay.Json{}.Msg(c, 200,
		"success", gin.H{
			"list": list,
		})

}
