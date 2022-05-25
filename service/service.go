package service

import (
	"github.com/gin-gonic/gin"
)

func Set(r *gin.Engine) *gin.Engine {
	r.Use(Cors())
	r.Use(Header())
	//r.Use(Pretreatment())
	return r
}
