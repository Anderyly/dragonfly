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

type CouponController struct {
}

type userCouponControllerForm struct {
	Page int `form:"page" binding:"required" label:"页码"`
}

func (con CouponController) User(c *gin.Context) {
	var data userCouponControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	page := data.Page - 1

	isAuth, code, msg, requestUser := Auth(GetToken(Token))

	if !isAuth {
		ay.Json{}.Msg(c, code, msg, gin.H{})
		return
	}

	var coupon []models.Coupon
	ay.Db.Select("id,name,sub_name,effective_at,des,status").
		Where("uid = ?", requestUser.Id).
		Offset(page * Limit).
		Limit(Limit).Debug().
		Order("status desc").
		Find(&coupon)

	var list []gin.H
	for _, v := range coupon {
		list = append(list, gin.H{
			"id":           v.Id,
			"name":         v.Name,
			"sub_name":     v.SubName,
			"des":          v.Des,
			"effective_at": v.EffectiveAt,
			"status":       v.Status,
		})
	}

	if list == nil {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"list": []string{},
		})
	} else {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"list": list,
		})
	}
}
