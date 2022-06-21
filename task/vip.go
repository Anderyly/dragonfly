/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package task

import (
	"dragonfly/ay"
	"dragonfly/controllers/api"
	"dragonfly/models"
	"log"
)

func Vip() {
	log.Println("vip开启每月发放")
	var res []models.User
	ay.Db.Where("now() < effective_at").Find(&res)
	//log.Println(res)
	for _, v := range res {
		api.CommonController{}.Give(v)
	}

}
