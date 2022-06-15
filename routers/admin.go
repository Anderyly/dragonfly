/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package routers

import (
	"dragonfly/controllers/admin"
	"github.com/gin-gonic/gin"
)

func AdminRouters(r *gin.RouterGroup) {

	adminGroup := r.Group("/admin/")
	// 登入
	adminGroup.POST("user/login", admin.Controller{}.Login)
	adminGroup.POST("pass", admin.Controller{}.Password)

	// 用户
	adminGroup.GET("account/list", admin.AccountController{}.List)
	adminGroup.POST("account/detail", admin.AccountController{}.Detail)
	adminGroup.POST("account/option", admin.AccountController{}.Option)

	// 次卡
	adminGroup.GET("card/list", admin.CardController{}.List)
	adminGroup.POST("card/detail", admin.CardController{}.Detail)
	adminGroup.POST("card/option", admin.CardController{}.Option)
	adminGroup.POST("card/delete", admin.CardController{}.Delete)

	// 等级
	adminGroup.GET("level/list", admin.LevelController{}.List)
	adminGroup.POST("level/detail", admin.LevelController{}.Detail)
	adminGroup.POST("level/option", admin.LevelController{}.Option)
	adminGroup.POST("level/delete", admin.LevelController{}.Delete)

	// vip
	adminGroup.GET("vip/list", admin.VipController{}.List)
	adminGroup.POST("vip/detail", admin.VipController{}.Detail)
	adminGroup.POST("vip/option", admin.VipController{}.Option)
	adminGroup.POST("vip/delete", admin.VipController{}.Delete)

	// Banner
	adminGroup.GET("banner/list", admin.BannerController{}.List)
	adminGroup.POST("banner/detail", admin.BannerController{}.Detail)
	adminGroup.POST("banner/option", admin.BannerController{}.Option)
	adminGroup.POST("banner/delete", admin.BannerController{}.Delete)
	adminGroup.POST("banner/upload", admin.BannerController{}.Upload)

	// 私教
	adminGroup.GET("education/list", admin.EducationController{}.List)
	adminGroup.GET("education/getAll", admin.EducationController{}.GetAll)
	adminGroup.POST("education/detail", admin.EducationController{}.Detail)
	adminGroup.POST("education/option", admin.EducationController{}.Option)
	adminGroup.POST("education/delete", admin.EducationController{}.Delete)
	adminGroup.POST("education/upload", admin.EducationController{}.Upload)

	// 拼课
	adminGroup.GET("spell/list", admin.SpellingController{}.List)
	adminGroup.POST("spell/detail", admin.SpellingController{}.Detail)
	adminGroup.POST("spell/option", admin.SpellingController{}.Option)
	adminGroup.POST("spell/delete", admin.SpellingController{}.Delete)
	adminGroup.POST("spell/upload", admin.SpellingController{}.Upload)

	// 拼课预约
	adminGroup.GET("spell_subscribe/list", admin.SpellingSubscribeController{}.List)

	// 拼课提问
	adminGroup.GET("spell_ask/list", admin.SpellingAskController{}.List)
	adminGroup.POST("spell_ask/detail", admin.SpellingAskController{}.Detail)
	adminGroup.POST("spell_ask/option", admin.SpellingAskController{}.Option)
	adminGroup.POST("spell_ask/delete", admin.SpellingAskController{}.Delete)
	adminGroup.POST("spell_ask/upload", admin.SpellingAskController{}.Upload)

	// 训练营
	adminGroup.GET("train/list", admin.TrainController{}.List)
	adminGroup.POST("train/detail", admin.TrainController{}.Detail)
	adminGroup.POST("train/option", admin.TrainController{}.Option)
	adminGroup.POST("train/delete", admin.TrainController{}.Delete)
	adminGroup.POST("train/upload", admin.TrainController{}.Upload)

	// 拼课预约
	adminGroup.GET("train_subscribe/list", admin.TrainSubscribeController{}.List)

	// 分店
	adminGroup.GET("store/getAll", admin.StoreController{}.GetAll)
	adminGroup.GET("store/list", admin.StoreController{}.List)
	adminGroup.POST("store/detail", admin.StoreController{}.Detail)
	adminGroup.POST("store/option", admin.StoreController{}.Option)
	adminGroup.POST("store/delete", admin.StoreController{}.Delete)
	adminGroup.POST("store/upload", admin.StoreController{}.Upload)

	// 舞蹈室
	adminGroup.GET("room/list", admin.RoomController{}.List)
	adminGroup.POST("room/detail", admin.RoomController{}.Detail)
	adminGroup.POST("room/option", admin.RoomController{}.Option)
	adminGroup.POST("room/delete", admin.RoomController{}.Delete)
	adminGroup.POST("room/upload", admin.RoomController{}.Upload)
	adminGroup.GET("room_subscribe/list", admin.RoomSubscribeController{}.List)

	adminGroup.POST("basic/detail", admin.BasicController{}.Detail)
	adminGroup.POST("basic/option", admin.BasicController{}.Option)

}
