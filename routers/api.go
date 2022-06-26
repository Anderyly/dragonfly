/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package routers

import (
	"dragonfly/controllers/api"
	"github.com/gin-gonic/gin"
)

func ApiRouters(r *gin.RouterGroup) {

	apiGroup := r.Group("/api/")

	apiGroup.POST("login", api.LoginController{}.Login)
	apiGroup.GET("main", api.HomeController{}.Main)
	apiGroup.GET("chart", api.HomeController{}.Chart)

	apiGroup.GET("level", api.VipController{}.Level)

	apiGroup.GET("test", api.HomeController{}.Test)

	apiGroup.GET("banner", api.HomeController{}.Banner)

	// 用户
	apiGroup.POST("user/info", api.UserController{}.Info)
	apiGroup.GET("user/order", api.UserController{}.Order)
	apiGroup.GET("user/qr", api.UserController{}.GetQr)
	apiGroup.GET("user/open", api.UserController{}.Open)

	// 分店
	apiGroup.GET("branch/list", api.StoreController{}.List)
	apiGroup.POST("branch/detail", api.StoreController{}.Detail)

	// 舞蹈室
	apiGroup.GET("room/list", api.RoomController{}.List)
	apiGroup.POST("room/detail", api.RoomController{}.Detail)
	apiGroup.GET("room/subscribe", api.RoomController{}.Subscribe)
	apiGroup.GET("room/subscribe_order", api.RoomController{}.SubscribeOrder)
	apiGroup.POST("room/pay", api.RoomController{}.Pay) // 测试
	apiGroup.POST("room/channel", api.RoomController{}.Channel)

	// 私教
	apiGroup.GET("education/list", api.EducationController{}.List)
	apiGroup.POST("education/detail", api.EducationController{}.Detail)

	// 拼课
	apiGroup.GET("class_spelling/list", api.SpellingController{}.List)
	apiGroup.POST("class_spelling/ask", api.SpellingController{}.Ask)
	apiGroup.POST("class_spelling/detail", api.SpellingController{}.Detail)
	apiGroup.POST("class_spelling/subscribe", api.SpellingController{}.Subscribe)

	// 训练营
	apiGroup.GET("train/list", api.TrainController{}.List)
	apiGroup.POST("train/detail", api.TrainController{}.Detail)
	apiGroup.POST("train/subscribe", api.TrainController{}.Subscribe)
	apiGroup.POST("train/channel", api.TrainController{}.Channel)

	// 反馈
	apiGroup.GET("feedback/list", api.FeedbackController{}.List)
	apiGroup.POST("feedback/add", api.FeedbackController{}.Add)

	// 优惠卷
	apiGroup.GET("coupon/user", api.CouponController{}.User)
	apiGroup.POST("coupon/select", api.CouponController{}.Select)

	// 次卡
	apiGroup.GET("card/list", api.CardController{}.List)
	apiGroup.GET("card/user", api.CardController{}.User)
	apiGroup.POST("card/pay", api.CardController{}.Pay)
	apiGroup.POST("card/select", api.CardController{}.Select)

	// vip
	apiGroup.GET("vip/amount", api.VipController{}.Amount)
	apiGroup.POST("vip/pay", api.VipController{}.Pay)

	// 充值
	apiGroup.GET("recharge/amount", api.RechargeController{}.Amount)
	apiGroup.POST("recharge/pay", api.RechargeController{}.Pay)

	// 异步
	apiGroup.Any("notify/wechat", api.NotifyController{}.Wechat)

}
