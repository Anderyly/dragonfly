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
	"strconv"
	"strings"
	"time"
)

type TrainController struct {
}

type listTrainControllerForm struct {
	StoreId int `form:"store_id" binding:"required" label:"分店id"`
	Page    int `form:"page" binding:"required" label:"页码"`
}

func (con TrainController) List(c *gin.Context) {
	var data listTrainControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	page := data.Page - 1

	var list []models.Train
	ay.Db.Select("id,name,start_date,back_image,amount,address,max_user").
		Where("store_id = ?", data.StoreId).
		Offset(page * Limit).
		Order("start_date desc").
		Limit(Limit).
		Find(&list)

	var r []gin.H

	_, _, _, requestUser := Auth(GetToken(Token))

	// 获取会员等级
	isVip, vipLevel := CommonController{}.GetVip(requestUser)

	for _, v := range list {
		vipAmount := v.Amount
		if isVip {
			vipAmount = v.Amount * (vipLevel.SpellClass / 10)
		}
		status, _, _ := con.GetStatus(requestUser, v)
		r = append(r, gin.H{
			"id":         v.Id,
			"name":       v.Name,
			"back_image": ay.Yaml.GetString("domain") + v.BackImage,
			"start_date": v.StartDate.Unix(),
			"amount":     v.Amount,
			"vip_amount": vipAmount,
			"status":     status,
			"address":    v.Address,
		})
	}

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"list": r,
	})

}

func (con TrainController) GetStatus(requestUser models.User, train models.Train) (int, int, int64) {
	status := 0
	userSubscribe := 0

	var cCount int64
	ay.Db.Model(&models.TrainLog{}).Where("train_id = ? AND uid = ?", train.Id, requestUser.Id).Count(&cCount)
	userSubscribe = int(cCount)
	if cCount > 0 {
		userSubscribe = 1
	}

	// 预约人数
	var count int64
	ay.Db.Model(&models.TrainLog{}).Where("train_id = ?", train.Id).Count(&count)
	log.Println(int(count))
	log.Println(train.MaxUser)
	if int(count) >= train.MaxUser {
		status = 2
	} else {
		status = 1
	}

	if time.Now().Unix() > train.StartDate.Unix() {
		status = 3
	}
	return status, userSubscribe, count
}

type detailTrainControllerForm struct {
	Id int `form:"id" binding:"required" label:"训练营id"`
}

func (con TrainController) Detail(c *gin.Context) {
	var data detailTrainControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	var res models.Train
	ay.Db.Where("id = ?", data.Id).First(&res)

	if res.Id == 0 {
		ay.Json{}.Msg(c, 400, "训练营不存在", gin.H{})
		return
	}

	_, _, _, requestUser := Auth(GetToken(Token))

	status, userSubscribe, count := con.GetStatus(requestUser, res)

	imageArr := strings.Split(res.Image, ",")

	var image []string

	for _, v := range imageArr {
		image = append(image, ay.Yaml.GetString("domain")+v)
	}

	video := ""
	if res.Video != "" {
		video = ay.Yaml.GetString("domain") + res.Video
	}

	var education models.Education
	ay.Db.First(&education, res.EducationId)

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"id":                  res.Id,
		"education_avatar":    ay.Yaml.GetString("domain") + education.Avatar,
		"education_name":      education.Name,
		"education_label":     education.Label,
		"education_introduce": education.Introduce,
		"name":                res.Name,
		"address":             res.Address,
		"mobile":              res.Mobile,
		"image":               image,
		"introduce":           res.Introduce,
		"attention":           res.Attention,
		"video":               video,
		"min_user":            res.MinUser,
		"max_user":            res.MaxUser,
		"room":                res.Room,
		"start_date":          res.StartDate.Unix(),
		"arrange":             res.Arrange,
		"user_subscribe":      userSubscribe,
		"status":              status,
		"count":               res.MaxUser - int(count),
		"amount":              res.Amount,
	})

}

type subscribeTrainControllerForm struct {
	Id       int    `form:"id" binding:"required" label:"训练营id"`
	Code     string `form:"code" binding:"required" label:"微信code"`
	PayType  int    `form:"pay_type" binding:"required" label:"支付方式"`
	CouponId int64  `form:"coupon_id" label:"优惠卷id"`
}

// Subscribe 预约、取消预约
func (con TrainController) Subscribe(c *gin.Context) {
	var data subscribeTrainControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	isLogin, code, msg, requestUser := Auth(GetToken(Token))

	if !isLogin {
		ay.Json{}.Msg(c, code, msg, gin.H{})
		return
	}

	var res models.Train
	ay.Db.Where("id = ?", data.Id).First(&res)

	if res.Id == 0 {
		ay.Json{}.Msg(c, 400, "训练营不存在", gin.H{})
		return
	}

	//// 获取会员等级
	isVip, vipLevel := CommonController{}.GetVip(requestUser)

	var trainLog models.TrainLog
	ay.Db.Where("train_id = ? AND uid = ?", data.Id, requestUser.Id).First(&trainLog)
	if trainLog.Id != 0 {
		ay.Json{}.Msg(c, 400, "您已预约", gin.H{})
		return
	}

	// 开始预约
	// 预约人数
	var count int64
	ay.Db.Model(&models.ClassSpellingLog{}).Where("spelling_id = ?", data.Id).Count(&count)
	if int(count) >= res.MaxUser {
		ay.Json{}.Msg(c, 400, "预约人数已满", gin.H{})
		return
	}

	// 优惠卷判断
	isCoupon, couponMsg, coupon := CommonController{}.IsCoupon(requestUser.Id, data.CouponId, 4)

	if data.CouponId != 0 {
		if !isCoupon {
			ay.Json{}.Msg(c, 400, couponMsg, gin.H{})
			return
		}
	}

	amount := res.Amount - coupon.Amount

	if isVip {
		amount = amount * (vipLevel.SpellClass / 10)
	}

	if amount <= 0.01 {
		amount = 0.01
	}
	isOrder, outTradeNo, order := CommonController{}.MakeOrder(2, 4, "训练营预约", requestUser.Id, coupon.Id, amount, res.Amount, c.ClientIP(), "", 0, res.Id, vipLevel.SpellClass)
	if isOrder == 0 {
		ay.Json{}.Msg(c, 400, outTradeNo, gin.H{})
		return
	}

	if data.PayType == 1 {
		isC, msg, info := PayController{}.Wechat(data.Code, order.Des, order.OutTradeNo, c.ClientIP(), amount)

		if !isC {
			ay.Json{}.Msg(c, 400, msg, gin.H{})
			return
		}

		ay.Json{}.Msg(c, 200, "success", gin.H{
			"info": info,
		})
	} else if data.PayType == 2 {
		if requestUser.Amount < amount {
			ay.Json{}.Msg(c, 400, "账户余额不足", gin.H{})
			return
		}
		requestUser.Amount -= amount
		requestUser.VipNum += res.Amount
		tx := ay.Db.Begin()

		// 扣除用户余额
		if err := tx.Save(&requestUser).Error; err == nil {
			if err := tx.Create(&models.TrainLog{
				TrainId: res.Id,
				Uid:     requestUser.Id,
				Oid:     order.Id,
			}).Error; err != nil {
				tx.Rollback()
				ay.Json{}.Msg(c, 400, "预约失败", gin.H{})
				return
			} else {
				if coupon.Id != 0 {
					coupon.Status = 1
					stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"), time.Local)
					coupon.UseAt = models.MyTime{Time: stamp}
				}
				if err := tx.Save(&coupon).Error; err != nil {
					tx.Rollback()
					ay.Json{}.Msg(c, 400, "请联系管理员", gin.H{})
					return
				} else {
					order.Status = 1
					stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"), time.Local)
					order.PayAt = models.MyTime{Time: stamp}
					if err := tx.Save(&order).Error; err != nil {
						tx.Rollback()
						ay.Json{}.Msg(c, 400, "请联系管理员", gin.H{})
						return
					} else {
						tx.Commit()
						ay.Json{}.Msg(c, 200, "已预约成功", gin.H{})
						return
					}

				}
			}
		} else {
			tx.Rollback()
			ay.Json{}.Msg(c, 400, "请联系管理员", gin.H{})
		}
	} else {
		ay.Json{}.Msg(c, 400, "支付方式不存在", gin.H{})
	}

}

type channelTrainControllerForm struct {
	Id int `form:"id" binding:"required" label:"训练营id"`
}

// Channel 取消预约
func (con TrainController) Channel(c *gin.Context) {
	var data channelTrainControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	isLogin, code, msg, requestUser := Auth(GetToken(Token))

	if !isLogin {
		ay.Json{}.Msg(c, code, msg, gin.H{})
		return
	}

	var res models.Train
	ay.Db.Where("id = ?", data.Id).First(&res)

	if res.Id == 0 {
		ay.Json{}.Msg(c, 400, "训练营不存在", gin.H{})
		return
	}

	//// 获取会员等级
	isVip, vipLevel := CommonController{}.GetVip(requestUser)

	var trainLog models.TrainLog
	ay.Db.Where("train_id = ? AND uid = ?", data.Id, requestUser.Id).First(&trainLog)

	if trainLog.Id == 0 {
		ay.Json{}.Msg(c, 400, "您还没有预约", gin.H{})
		return
	}

	var order models.Order
	ay.Db.First(&order, trainLog.Oid)

	if order.Id == 0 || order.Status != 1 {
		ay.Json{}.Msg(c, 400, "您还没有预约", gin.H{})
		return
	}

	// 取消预约
	channelNum := ay.Yaml.GetInt("channel.num")
	channelHour := ay.Yaml.GetInt("channel.hour")
	if isVip {
		channelNum = vipLevel.Chanel
		channelHour = vipLevel.AdvanceChanel
	}

	if res.StartDate.Unix() < time.Now().Unix() {
		ay.Json{}.Msg(c, 400, "训练营已开始，无法取消", gin.H{})
		return
	}

	log.Println(res.StartDate.Unix() - time.Now().Unix())
	log.Println(channelHour * 3600)
	if res.StartDate.Unix()-time.Now().Unix() < int64(channelHour*3600) {
		ay.Json{}.Msg(c, 400, "距离开始不足"+strconv.Itoa(channelHour)+"小时，无法取消", gin.H{})
		return
	}

	// 是否还能取消
	isChannel, channelMsg := CommonController{}.GetChanelCount(channelNum, requestUser.Id)

	if !isChannel {
		ay.Json{}.Msg(c, 400, channelMsg, gin.H{})
		return
	}

	tx := ay.Db.Begin()

	// 查询当前用户取消预约次数
	if err := tx.Delete(&trainLog).Error; err != nil {
		tx.Rollback()
		ay.Json{}.Msg(c, 400, "取消预约失败", gin.H{})
	} else {
		if err := tx.Create(&models.UserChannel{
			Type: 1,
			Cid:  res.Id,
			Uid:  requestUser.Id,
		}).Error; err != nil {
			tx.Rollback()
		} else {
			is, msg := CommonController{}.Back(order.Id)
			if !is {
				tx.Rollback()
				ay.Json{}.Msg(c, 200, msg, gin.H{})
			} else {
				tx.Commit()
				ay.Json{}.Msg(c, 200, "已取消预约", gin.H{})
			}

		}
	}

}
