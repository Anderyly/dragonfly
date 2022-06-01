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
)

type FeedbackController struct {
}

type listFeedbackControllerForm struct {
	Page int `form:"page" binding:"required" label:"页码"`
}

func (con FeedbackController) List(c *gin.Context) {
	var data listFeedbackControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	isAuth, code, msg, requestUser := Auth(GetToken(Token))

	if !isAuth {
		ay.Json{}.Msg(c, code, msg, gin.H{})
		return
	}

	page := data.Page - 1

	var list []models.Feedback
	ay.Db.Where("uid = ?", requestUser.Id).
		Offset(page * Limit).
		Limit(Limit).Debug().
		Order("created_at desc").
		Find(&list)

	var r []gin.H

	for _, v := range list {

		r = append(r, gin.H{
			"id":         v.Id,
			"name":       v.Name,
			"mobile":     v.Mobile,
			"status":     v.Status,
			"created_at": v.CreatedAt,
		})
	}

	if r == nil {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"list": []string{},
		})
	} else {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"list": r,
		})
	}

}

type addFeedbackControllerForm struct {
	Name    string `form:"name" binding:"required" label:"姓名"`
	Mobile  string `form:"mobile" binding:"required" label:"手机"`
	Content string `form:"content" binding:"required" label:"内容"`
}

func (con FeedbackController) Add(c *gin.Context) {
	var data addFeedbackControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	isAuth, code, msg, requestUser := Auth(GetToken(Token))

	if !isAuth {
		ay.Json{}.Msg(c, code, msg, gin.H{})
		return
	}

	if err := ay.Db.Create(&models.Feedback{
		Uid:     requestUser.Id,
		Name:    data.Name,
		Mobile:  data.Mobile,
		Content: data.Content,
		Status:  0,
	}).Error; err != nil {
		ay.Json{}.Msg(c, 400, "反馈失败", gin.H{})
	} else {
		ay.Json{}.Msg(c, 200, "反馈成功", gin.H{})
	}

}
