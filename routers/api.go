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
	apiGroup.POST("branch/list", api.StoreController{}.List)
	apiGroup.POST("branch/detail", api.StoreController{}.Detail)
}
