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

type FeedController struct {
}

// List 列表
func (con FeedController) List(c *gin.Context) {

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

	var list []models.Feedback

	var count int64
	res := ay.Db
	//if data.Key != "" {
	//	res.Where("name like ?", "%"+data.Key+"%")
	//}

	row := res

	row.Count(&count)

	res.Order("created_at asc").
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
func (con FeedController) Detail(c *gin.Context) {
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

	var user models.Feedback

	ay.Db.First(&user, data.Id)

	ay.Json{}.Msg(c, 200,
		"success", gin.H{
			"info": user,
		})
}

// Option 添加 编辑
func (con FeedController) Option(c *gin.Context) {

	type optionForm struct {
		Id        int64   `form:"id"`
		Name      string  `form:"name"`
		Amount    float64 `form:"amount"`
		OldAmount float64 `form:"old_amount"`
		Discount  float64 `form:"discount"`
		Sort      int     `form:"sort"`
		Effective int     `form:"effective"`
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

	var res models.Vip
	ay.Db.First(&res, data.Id)
	if res.Id != 0 {
		res.Name = data.Name
		res.Amount = data.Amount
		res.OldAmount = data.OldAmount
		res.Discount = data.Discount
		res.Sort = data.Sort
		res.Effective = data.Effective

		ay.Db.Save(&res)
		ay.Json{}.Msg(c, 200, "修改成功", gin.H{})
	} else {
		ay.Db.Create(&models.Vip{
			Name:      data.Name,
			Amount:    data.Amount,
			OldAmount: data.OldAmount,
			Discount:  data.Discount,
			Sort:      data.Sort,
			Effective: data.Effective,
		})
		ay.Json{}.Msg(c, 200, "创建成功", gin.H{})
	}

}

// Delete 删除
func (con FeedController) Delete(c *gin.Context) {
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
		var res models.Vip
		ay.Db.Delete(&res, v)
	}

	ay.Json{}.Msg(c, 200,
		"删除成功", gin.H{})
}
