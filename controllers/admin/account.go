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

type AccountController struct {
}

// List 列表
func (con AccountController) List(c *gin.Context) {

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

	var user []models.User

	var count int64
	res := ay.Db.Model(models.User{})
	if data.Key != "" {
		res.Where("nickname like ?", "%"+data.Key+"%")
	}

	row := res

	row.Count(&count)

	res.Order("created_at desc").
		Limit(data.PageSize).
		Offset((data.Page - 1) * data.PageSize).
		Find(&user)

	var list []gin.H

	for _, v := range user {
		var vipLevel models.VipLevel

		ay.Db.Order("num desc").
			Select("name").
			Where("num <= ?", v.VipNum).
			First(&vipLevel)

		list = append(list, gin.H{
			"id":           v.Id,
			"vip_num":      v.VipNum,
			"amount":       v.Amount,
			"avatar":       v.Avatar,
			"nickname":     v.NickName,
			"level":        vipLevel.Name,
			"created_at":   v.CreatedAt,
			"effective_at": v.EffectiveAt,
		})
	}

	ay.Json{}.Msg(c, 200,
		"success", gin.H{
			"list":  list,
			"total": count,
		})
}

// Detail 详情
func (con AccountController) Detail(c *gin.Context) {
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

	var user models.User

	ay.Db.First(&user, data.Id)

	ay.Json{}.Msg(c, 200,
		"success", gin.H{
			"info": user,
		})
}

// Option 添加 编辑
func (con AccountController) Option(c *gin.Context) {

	type optionForm struct {
		Id          int64   `form:"id"`
		Avatar      string  `form:"avatar"`
		Nickname    string  `form:"nickname"`
		Amount      float64 `form:"amount"`
		VipNum      float64 `form:"vip_num"`
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

	var user models.User
	ay.Db.First(&user, data.Id)

	stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", data.EffectiveAt, time.Local)

	user.EffectiveAt = models.MyTime{Time: stamp}

	user.Avatar = data.Avatar
	user.NickName = data.Nickname
	user.Amount = data.Amount
	user.VipNum = data.VipNum

	ay.Db.Save(&user)
	ay.Json{}.Msg(c, 200, "修改成功", gin.H{})

}

// Delete 删除
func (con AccountController) Delete(c *gin.Context) {
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
		var user models.User
		ay.Db.Delete(&user, v)
	}

	ay.Json{}.Msg(c, 200,
		"删除成功", gin.H{})
}
