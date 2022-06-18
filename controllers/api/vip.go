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
	"log"
	"time"
)

type VipController struct {
}

func (con VipController) Amount(c *gin.Context) {

	var vip []models.Vip

	ay.Db.Order("sort asc").Find(&vip)

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"list": vip,
	})
}

type payVipControllerForm struct {
	Id      int64  `form:"id" binding:"required" label:"id"`
	PayType int    `form:"pay_type" label:"支付方式"`
	Code    string `form:"code" binding:"required" label:"微信code"`
}

func (con VipController) Pay(c *gin.Context) {
	var data payVipControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	isAuth, code, msg, requestUser := Auth(GetToken(Token))

	if !isAuth {
		ay.Json{}.Msg(c, code, msg, gin.H{})
		return
	}

	var vip models.Vip
	ay.Db.First(&vip, data.Id)
	if vip.Id == 0 {
		ay.Json{}.Msg(c, 400, "支付id错误", gin.H{})
		return
	}

	isOrder, outTradeNo, order := CommonController{}.MakeOrder(2, 3, "开通"+vip.Name, requestUser.Id, 0, vip.Amount, vip.Amount, c.ClientIP(), "", 0, vip.Id, 0)
	if isOrder == 0 {
		ay.Json{}.Msg(c, 400, outTradeNo, gin.H{})
		return
	}

	if data.PayType == 1 {
		isC, msg, info := PayController{}.Wechat(data.Code, order.Des, order.OutTradeNo, c.ClientIP(), vip.Amount)

		if !isC {
			ay.Json{}.Msg(c, 400, msg, gin.H{})
			return
		}

		ay.Json{}.Msg(c, 200, "success", gin.H{
			"info": info,
		})
	} else {
		is, msg, user := CommonController{}.UserAmountPay(&requestUser, vip.Amount)
		log.Println(user)
		if is {
			tx := ay.Db.Begin()
			order.Status = 1
			stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"), time.Local)
			order.PayAt = models.MyTime{Time: stamp}
			if err := tx.Save(&order).Error; err != nil {
				tx.Rollback()
				ay.Json{}.Msg(c, 400, "请联系管理员", gin.H{})
				return
			}
			isUserVip := models.UserModel{}.SetVip(requestUser, vip.Effective, vip.Amount)
			if !isUserVip {
				tx.Rollback()
				ay.Json{}.Msg(c, 400, "请联系管理员", gin.H{})
				return
			}
			tx.Commit()

			var user1 models.User
			ay.Db.First(&user1, requestUser.Id)
			CommonController{}.Give(user1)

			ay.Json{}.Msg(c, 200, "用户支付成功", gin.H{})
		} else {
			ay.Json{}.Msg(c, 400, msg, gin.H{})
		}

	}

}
