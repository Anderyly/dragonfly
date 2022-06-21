package main

import (
	"dragonfly/ay"
	"dragonfly/routers"
	"dragonfly/service"
	"dragonfly/task"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
)

var (
	r *gin.Engine
)

func main() {
	//service.ElectricServer{}.Header("/cloud/thing/properties/get", "")
	//return
	ay.Yaml = ay.InitConfig()
	ay.Sql()
	go ay.WatchConf()

	// 开启定时任务
	c := cron.New()
	_, err := c.AddFunc("@every 5m", task.Start)
	if err != nil {
		log.Println(err.Error())
		return
	}
	c.Start()

	r = gin.Default()
	gin.SetMode(gin.DebugMode)
	r = service.Set(r)
	r.StaticFS("/static/", http.Dir("./static"))
	r = routers.GinRouter(r)

	err = r.Run(":" + ay.Yaml.GetString("port"))
	if err != nil {
		panic(err.Error())
	}
	select {}
}
