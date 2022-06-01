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

	apiGroup.GET("banner", api.HomeController{}.Banner)

	// 用户
	apiGroup.POST("user/info", api.UserController{}.Info)

	// 分店
	apiGroup.GET("branch/list", api.StoreController{}.List)
	apiGroup.POST("branch/detail", api.StoreController{}.Detail)

	// 舞蹈室
	apiGroup.GET("room/list", api.RoomController{}.List)
	apiGroup.POST("room/detail", api.RoomController{}.Detail)

	// 私教
	apiGroup.GET("education/list", api.EducationController{}.List)
	apiGroup.POST("education/detail", api.EducationController{}.Detail)

	// 拼课
	apiGroup.GET("class_spelling/list", api.SpellingController{}.List)
	apiGroup.POST("class_spelling/ask", api.SpellingController{}.Ask)
	apiGroup.POST("class_spelling/detail", api.SpellingController{}.Detail)
	apiGroup.POST("class_spelling/subscribe", api.SpellingController{}.Subscribe)

	// 反馈
	apiGroup.GET("feedback/list", api.FeedbackController{}.List)
	apiGroup.POST("feedback/add", api.FeedbackController{}.Add)

	// 优惠卷
	apiGroup.GET("coupon/user", api.CouponController{}.User)

	// 次卡
	apiGroup.GET("card/list", api.CardController{}.List)
	apiGroup.GET("card/user", api.CardController{}.User)
}
