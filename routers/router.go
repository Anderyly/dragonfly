package routers

import (
	"github.com/gin-gonic/gin"
)

func GinRouter(r *gin.Engine) *gin.Engine {

	router := r.Group("/")

	ApiRouters(router)
	AdminRouters(router)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "访问错误",
		})
	})

	return r
}
