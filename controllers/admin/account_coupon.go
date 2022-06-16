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
	"strconv"
	"strings"
	"time"
)

type AccountCouponController struct {
}

// List 列表
func (con AccountCouponController) List(c *gin.Context) {

	type listForm struct {
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
		Key      string `form:"key"`
		Uid      int64  `form:"uid"`
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
		models.Coupon
		TypeName string `json:"type_name"`
	}
	var list []r

	var count int64
	res := ay.Db.Table("d_user_coupon").Select("d_user_coupon.*,d_user_coupon.type as type_name")

	res.Where("d_user_coupon.uid = ?", data.Uid)
	row := res

	row.Count(&count)

	res.Order("d_user_coupon.created_at desc").
		Limit(data.PageSize).
		Offset((data.Page - 1) * data.PageSize).
		Find(&list)

	for k, v := range list {
		typeN := []string{
			"教室", "次卡", "会员", "训练营", "充值",
		}
		var typeName string

		for _, v1 := range strings.Split(v.Type, ",") {
			s, _ := strconv.Atoi(v1)
			typeName += typeN[s-1] + "、"
		}
		list[k].TypeName = strings.TrimRight(typeName, "、")
	}

	ay.Json{}.Msg(c, 200,
		"success", gin.H{
			"list":  list,
			"total": count,
		})
}

// Detail 详情
func (con AccountCouponController) Detail(c *gin.Context) {
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

	var user models.Coupon

	ay.Db.First(&user, data.Id)

	ay.Json{}.Msg(c, 200,
		"success", gin.H{
			"info": user,
		})
}

// Option 添加 编辑
func (con AccountCouponController) Option(c *gin.Context) {

	type optionForm struct {
		Id          int64   `form:"id"`
		Type        string  `form:"type"`
		Uid         int64   `form:"uid"`
		Name        string  `form:"name"`
		SubName     string  `form:"sub_name"`
		Des         string  `form:"des"`
		Amount      float64 `form:"amount"`
		Status      int     `form:"status"`
		EffectiveAt string  `form:"effective_at"`
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

	stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", data.EffectiveAt, time.Local)

	var res models.Coupon
	ay.Db.First(&res, data.Id)
	if res.Id != 0 {
		res.EffectiveAt = models.MyTime{Time: stamp}
		res.Des = data.Des
		res.Status = data.Status
		res.Type = data.Type
		res.Name = data.Name
		res.SubName = data.SubName
		res.Amount = data.Amount

		ay.Db.Save(&res)
		ay.Json{}.Msg(c, 200, "修改成功", gin.H{})
	} else {
		ay.Db.Create(&models.Coupon{
			Type:        data.Type,
			Uid:         data.Uid,
			Name:        data.Name,
			SubName:     data.SubName,
			Des:         data.Des,
			Amount:      data.Amount,
			EffectiveAt: models.MyTime{Time: stamp},
			Status:      data.Status,
		})
		ay.Json{}.Msg(c, 200, "创建成功", gin.H{})

	}

}

// Delete 删除
func (con AccountCouponController) Delete(c *gin.Context) {
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
		var user models.Coupon
		ay.Db.Delete(&user, v)
	}

	ay.Json{}.Msg(c, 200,
		"删除成功", gin.H{})
}
