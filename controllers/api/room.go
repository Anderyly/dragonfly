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
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"strings"
	"time"
)

type RoomController struct {
}

type listRoomControllerForm struct {
	StoreId int `form:"store_id" binding:"required" label:"分店id"`
	Page    int `form:"page" binding:"required" label:"页码"`
}

func (con RoomController) List(c *gin.Context) {
	var data listRoomControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	page := data.Page - 1

	var list []models.Room
	ay.Db.Select("id,name,size,image,address").
		Where("store_id = ?", data.StoreId).
		Offset(page * Limit).
		Limit(Limit).
		Find(&list)

	var r []gin.H

	for _, v := range list {
		imageArr := strings.Split(v.Image, ",")

		r = append(r, gin.H{
			"id":      v.Id,
			"name":    v.Name,
			"size":    v.Size,
			"label":   "火爆",
			"image":   ay.Yaml.GetString("domain") + imageArr[0],
			"address": v.Address,
		})
	}

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"list": r,
	})

}

type detailRoomControllerForm struct {
	Id int `form:"id" binding:"required" label:"舞蹈室id"`
}

func (con RoomController) Detail(c *gin.Context) {
	var data detailRoomControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	var res models.Room
	ay.Db.Where("id = ?", data.Id).First(&res)

	if res.Id == 0 {
		ay.Json{}.Msg(c, 400, "舞蹈室不存在", gin.H{})
		return
	}

	imageArr := strings.Split(res.Image, ",")

	var image []string

	for _, v := range imageArr {
		image = append(image, ay.Yaml.GetString("domain")+v)
	}

	video := ""
	if res.Video != "" {
		video = ay.Yaml.GetString("domain") + res.Video
	}

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"id":        res.Id,
		"name":      res.Name,
		"address":   res.Address,
		"mobile":    res.Mobile,
		"image":     image,
		"introduce": res.Introduce,
		"attention": res.Attention,
		"guide":     res.Guide,
		"cancel":    res.Cancel,
		"video":     video,
		"type":      res.Type,
	})

}

type subscribeRoomControllerForm struct {
	Id  int64 `form:"id" binding:"required" label:"舞蹈室id"`
	Ymd int   `form:"ymd" binding:"required" label:"时间"`
}

func (con RoomController) Subscribe(c *gin.Context) {
	var data subscribeRoomControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	var res models.Room
	ay.Db.Where("id = ?", data.Id).First(&res)

	if res.Id == 0 {
		ay.Json{}.Msg(c, 400, "舞蹈室不存在", gin.H{})
		return
	}

	isLogin, code, msg, requestUser := Auth(GetToken(Token))

	if !isLogin {
		ay.Json{}.Msg(c, code, msg, gin.H{})
		return
	}

	// 查询有无时间
	hourId := con.getHour(data.Id, data.Ymd)

	var userSub []models.UserRoomSubscribe
	ay.Db.Where("uid = ? and ymd = ?", requestUser.Id, data.Ymd).Find(&userSub)

	//config := models.ConfigModel{}.Get()
	var list []gin.H
	for _, v := range userSub {
		amount := res.Amount
		var roomSub models.RoomSubscribe
		ay.Db.Where("room_id = ? and ymd = ? and hour_id = ?", data.Id, data.Ymd, v.HourId).Find(&roomSub)
		if roomSub.Id != 0 {
			amount = roomSub.Amount
		}
		//list[v.HourId] = v.Amount
		list = append(list, gin.H{
			"hour_id": v.HourId,
			"amount":  amount,
		})
	}

	if list != nil {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"hour":           hourId,
			"user_subscribe": list,
		})
	} else {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"hour":           hourId,
			"user_subscribe": gin.H{},
		})
	}

}

func (con RoomController) getHour(roomId int64, ymd int) []gin.H {
	var res []gin.H

	var room models.Room
	ay.Db.First(&room, roomId)
	//config := models.ConfigModel{}.Get()
	// 后台定义价格
	var roomSub []models.RoomSubscribe
	ay.Db.Where("room_id = ? and ymd = ?", roomId, ymd).Find(&roomSub)

	var userSub []models.UserRoomSubscribe
	ay.Db.Where("room_id = ? and ymd = ?", roomId, ymd).Find(&userSub)

	for i := 1; i <= 24; i++ {
		amount := room.Amount
		isSubscribe := 0
		for _, v := range roomSub {
			if v.HourId == i {
				amount = v.Amount
			}
		}

		for _, v := range userSub {
			if v.HourId == i {
				isSubscribe = 1
			}
		}

		res = append(res, gin.H{
			"hour_id":      i,
			"amount":       amount,
			"is_subscribe": isSubscribe,
		})
	}

	return res
}

type payRoomControllerForm struct {
	RoomId   int64  `form:"room_id" binding:"required" label:"舞蹈室id"`
	HourId   string `form:"hour_id" binding:"required" label:"时间id"`
	Code     string `form:"code" binding:"required" label:"微信code"`
	PayType  int    `form:"pay_type" binding:"required" label:"支付方式"`
	CouponId int64  `form:"coupon_id" label:"优惠卷id"`
	CardId   int64  `form:"card_id" label:"次卡id"`
	Ymd      int    `form:"ymd" binding:"required" label:"日期"`
}

func (con RoomController) Pay(c *gin.Context) {
	var data payRoomControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	var res models.Room
	ay.Db.Where("id = ?", data.RoomId).First(&res)

	if res.Id == 0 {
		ay.Json{}.Msg(c, 400, "舞蹈室不存在", gin.H{})
		return
	}

	isLogin, code, msg, requestUser := Auth(GetToken(Token))

	if !isLogin {
		ay.Json{}.Msg(c, code, msg, gin.H{})
		return
	}

	amount := 0.00

	hourIdArr := strings.Split(data.HourId, ",")
	oldHourIdArr := hourIdArr

	log.Println(oldHourIdArr)

	// 判断是否有用户在时间预约
	var userSub []models.UserRoomSubscribe
	ay.Db.Where("room_id = ? and ymd = ?", data.RoomId, data.Ymd).Find(&userSub)
	for _, v := range userSub {
		for _, v1 := range hourIdArr {
			if strconv.Itoa(v.HourId) == v1 {
				ay.Json{}.Msg(c, 400, "当前时间段已有预约", gin.H{})
				return
			}
		}
	}
	num := 0

	var userCard models.UserCard

	// 去除次卡价格
	if data.CardId != 0 {
		ay.Db.First(&userCard, data.CardId)
		if userCard.Status == 1 {
			ay.Json{}.Msg(c, 400, "次卡已用完", gin.H{})
			return
		}
		if userCard.Status == 2 {
			ay.Json{}.Msg(c, 400, "次卡已过期", gin.H{})
			return
		}
		if userCard.Status == 2 {
			ay.Json{}.Msg(c, 400, "次卡已过期", gin.H{})
			return
		}
		var card models.Card
		ay.Db.First(&card, userCard.CardId)
		if card.Type != 5 {
			if card.Type != res.Type {
				ay.Json{}.Msg(c, 400, "次卡类型不对", gin.H{})
				return
			}
		}
		cardNum := userCard.AllHour - userCard.UseHour
		log.Println(cardNum)

		if cardNum >= len(hourIdArr) {
			num = len(hourIdArr)
			hourIdArr = []string{}
		} else {
			num = cardNum
			hourIdArr = hourIdArr[0 : len(hourIdArr)-cardNum]
		}

	}

	log.Println(hourIdArr)

	ssv := gin.H{
		"time": oldHourIdArr,
		"num":  num,
		"ymd":  data.Ymd,
	}

	isCoupon, couponMsg, coupon := CommonController{}.IsCoupon(requestUser.Id, data.CouponId, 1)
	isVip, vipLevel := CommonController{}.GetVip(requestUser)

	if len(hourIdArr) >= 1 {
		// 获取预定价格
		//var roomSub []models.RoomSubscribe
		//ay.Db.Where("room_id = ? and ymd = ?", data.RoomId, data.Ymd).Find(&roomSub)
		//for _, v := range roomSub {
		for _, v1 := range hourIdArr {
			var roomSub models.RoomSubscribe
			ay.Db.Where("room_id = ? and ymd = ? and hour_id = ?", data.RoomId, data.Ymd, v1).First(&roomSub)
			if roomSub.Id != 0 {
				//	log.Println(amount)
				amount += roomSub.Amount
			} else {
				amount += res.Amount
			}
		}
		//}

		log.Println("old:", amount)
		// 优惠卷判断

		if data.CouponId != 0 {
			if !isCoupon {
				ay.Json{}.Msg(c, 400, couponMsg, gin.H{})
				return
			}
		}

		amount = amount - coupon.Amount

		//log.Println("coupon:", amount)
		// 会员折扣
		// 获取会员等级
		if isVip {
			amount = amount * (vipLevel.Subscribe / 10)
		}
	}

	log.Println(amount)
	//return

	amount, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", amount), 64)
	//log.Println("price:", amount)
	//return
	//log.Println("card:", num)

	cv, _ := json.Marshal(ssv)
	isOrder, outTradeNo, order := CommonController{}.MakeOrder(2, 1, "舞蹈室预约", requestUser.Id, coupon.Id, amount, amount, c.ClientIP(), string(cv), data.CardId, res.Id, vipLevel.SpellClass)
	if isOrder == 0 {
		ay.Json{}.Msg(c, 400, outTradeNo, gin.H{})
		return
	}

	stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"), time.Local)
	if amount == 0 {
		tx := ay.Db.Begin()

		// 优惠卷设置过期
		if coupon.Id != 0 {
			coupon.Status = 1
			coupon.UseAt = models.MyTime{Time: stamp}
			if err := tx.Save(&coupon).Error; err != nil {
				tx.Rollback()
				ay.Json{}.Msg(c, 400, "请联系管理员", gin.H{})
				return
			}
		}
		// 扣除次卡
		if userCard.Id != 0 {
			userCard.UseHour += num
			if userCard.UseHour == userCard.AllHour {
				userCard.Status = 2
			}
			if err := tx.Save(&userCard).Error; err != nil {
				tx.Rollback()
				ay.Json{}.Msg(c, 400, "请联系管理员", gin.H{})
				return
			}
		}
		// 增加预约次数
		isSub, _ := models.UserRoomSubscribeModel{}.Add(outTradeNo, string(cv), requestUser.Id, res.Id, data.Ymd)
		if !isSub {
			tx.Rollback()
			ay.Json{}.Msg(c, 400, "请联系管理员", gin.H{})
			return
		}

		order.PayAt = models.MyTime{Time: stamp}
		order.Status = 1
		tx.Save(&order)

		tx.Commit()
		ay.Json{}.Msg(c, 200, "预约成功", gin.H{})
		return
	}

	if data.PayType == 1 {

		log.Println(amount)
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
		if isVip {
			requestUser.VipNum += amount
		}
		tx := ay.Db.Begin()

		// 扣除用户余额
		if err := tx.Save(&requestUser).Error; err == nil {

			// 优惠卷设置过期
			if coupon.Id != 0 {
				coupon.Status = 1
				coupon.UseAt = models.MyTime{Time: stamp}
				if err := tx.Save(&coupon).Error; err != nil {
					tx.Rollback()
					ay.Json{}.Msg(c, 400, "请联系管理员", gin.H{})
					return
				}
			}

			// 扣除次卡
			if userCard.Id != 0 {
				userCard.UseHour += num
				if userCard.UseHour == userCard.AllHour {
					userCard.Status = 2
				}
				if err := tx.Save(&userCard).Error; err != nil {
					tx.Rollback()
					ay.Json{}.Msg(c, 400, "请联系管理员", gin.H{})
					return
				}
			}
			// 增加预约次数
			isSub, _ := models.UserRoomSubscribeModel{}.Add(outTradeNo, string(cv), requestUser.Id, res.Id, data.Ymd)
			if !isSub {
				tx.Rollback()
				ay.Json{}.Msg(c, 400, "请联系管理员", gin.H{})
				return
			}

			order.Status = 1
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

		} else {
			tx.Rollback()
			ay.Json{}.Msg(c, 400, "请联系管理员", gin.H{})
		}
	} else {
		ay.Json{}.Msg(c, 400, "支付方式不存在", gin.H{})
	}

}

type channelRoomControllerForm struct {
	Oid string `form:"oid" binding:"required" label:"订单号"`
}

// Channel 取消预定
func (con RoomController) Channel(c *gin.Context) {
	var data channelRoomControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	isAuth, code, msg, requestUser := Auth(GetToken(Token))

	if !isAuth {
		ay.Json{}.Msg(c, code, msg, gin.H{})
		return
	}

	// 订单判断处理
	var order models.Order
	ay.Db.Where("id = ?", data.Oid).First(&order)
	if order.Id == 0 {
		ay.Json{}.Msg(c, 400, "订单不存在", gin.H{})
		return
	}
	if order.Status != 1 {
		ay.Json{}.Msg(c, 400, "订单未支付或已退款", gin.H{})
		return
	}
	if order.Uid != requestUser.Id {
		ay.Json{}.Msg(c, 400, "不是您的订单", gin.H{})
		return
	}

	var res models.Room
	ay.Db.First(&res, order.Cid)
	//////

	//// 获取会员等级
	isVip, vipLevel := CommonController{}.GetVip(requestUser)

	// 取消预约
	channelNum := ay.Yaml.GetInt("channel.num")
	channelHour := ay.Yaml.GetInt("channel.hour")
	if isVip {
		channelNum = vipLevel.Chanel
		channelHour = vipLevel.AdvanceChanel
	}

	// 舞蹈室
	type cv struct {
		Num  int      `json:"num"`
		Time []string `json:"time"`
		Ymd  int      `json:"ymd"`
	}
	var cc cv
	json.Unmarshal([]byte(order.Content), &cc)

	hour := cc.Time[0]
	if len(cc.Time[0]) == 1 {
		hour = "0" + cc.Time[0]
	}

	ymd := strconv.Itoa(cc.Ymd)
	ymd = ymd[0:4] + "-" + ymd[4:6] + "-" + ymd[6:8]

	stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", ymd+" "+hour+":00:00", time.Local)

	log.Println(ymd + " " + hour + ":00:00")

	//if res.StartDate.Unix() < time.Now().Unix() {
	//	ay.Json{}.Msg(c, 400, "训练营已开始，无法取消", gin.H{})
	//	return
	//}
	//
	if stamp.Unix()-time.Now().Unix() < int64(channelHour*3600) {
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

	if err := tx.Create(&models.UserChannel{
		Type: 2,
		Cid:  res.Id,
		Uid:  requestUser.Id,
		Oid:  order.Id,
	}).Error; err != nil {
		tx.Rollback()
	} else {
		is, msg := CommonController{}.Back(order.Id)
		if !is {
			tx.Rollback()
			ay.Json{}.Msg(c, 200, msg, gin.H{})
		} else {
			// 删除预约次数
			isSub, _ := models.UserRoomSubscribeModel{}.Del(order.Content, requestUser.Id, order.Cid)
			if !isSub {
				tx.Rollback()
				ay.Json{}.Msg(c, 400, "请联系管理员", gin.H{})
				return
			}
			tx.Commit()
			ay.Json{}.Msg(c, 200, "已取消预约", gin.H{})
		}

	}

}

type subscribeOrderRoomControllerForm struct {
	Oid string `form:"oid" binding:"required" label:"订单号"`
}

func (con RoomController) SubscribeOrder(c *gin.Context) {
	var data subscribeOrderRoomControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	var order models.Order
	ay.Db.First(&order, data.Oid)
	if order.Id == 0 {
		ay.Json{}.Msg(c, 400, "订单不存在", gin.H{})
		return
	}

	var res models.Room
	ay.Db.Where("id = ?", order.Cid).First(&res)

	if res.Id == 0 {
		ay.Json{}.Msg(c, 400, "舞蹈室不存在", gin.H{})
		return
	}

	isLogin, code, msg, _ := Auth(GetToken(Token))

	if !isLogin {
		ay.Json{}.Msg(c, code, msg, gin.H{})
		return
	}

	type cv struct {
		Num  int      `json:"num"`
		Time []string `json:"time"`
		Ymd  int      `json:"ymd"`
	}
	var cc cv
	json.Unmarshal([]byte(order.Content), &cc)

	var userSub []models.UserRoomSubscribe
	ay.Db.Debug().Where("out_trade_no = ?", order.OutTradeNo).Find(&userSub)

	//config := models.ConfigModel{}.Get()
	var list []gin.H
	for _, v := range userSub {
		amount := res.Amount
		var roomSub models.RoomSubscribe
		ay.Db.Where("room_id = ? and ymd = ? and hour_id = ?", res.Id, cc.Ymd, v.HourId).Find(&roomSub)
		if roomSub.Id != 0 {
			amount = roomSub.Amount
		}
		//list[v.HourId] = v.Amount
		list = append(list, gin.H{
			"hour_id": v.HourId,
			"amount":  amount,
		})
	}

	imageArr := strings.Split(res.Image, ",")

	var image []string

	for _, v := range imageArr {
		image = append(image, ay.Yaml.GetString("domain")+v)
	}

	video := ""
	if res.Video != "" {
		video = ay.Yaml.GetString("domain") + res.Video
	}

	if list != nil {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"info": gin.H{
				"id":        res.Id,
				"name":      res.Name,
				"address":   res.Address,
				"mobile":    res.Mobile,
				"image":     image,
				"introduce": res.Introduce,
				"attention": res.Attention,
				"guide":     res.Guide,
				"cancel":    res.Cancel,
				"video":     video,
				"type":      res.Type,
			},
			"subscribe_ymd":  cc.Ymd,
			"user_subscribe": list,
		})
	} else {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"info": gin.H{
				"id":        res.Id,
				"name":      res.Name,
				"address":   res.Address,
				"mobile":    res.Mobile,
				"image":     image,
				"introduce": res.Introduce,
				"attention": res.Attention,
				"guide":     res.Guide,
				"cancel":    res.Cancel,
				"video":     video,
				"type":      res.Type,
			},
			"user_subscribe": gin.H{},
		})
	}
}
