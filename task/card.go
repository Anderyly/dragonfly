/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package task

import (
	"dragonfly/ay"
	"dragonfly/models"
	"log"
)

func Card() {
	log.Println("开始次卡设置过期")
	var res []models.UserCard
	ay.Db.Where("status = 0 AND now() > effective_at").Select("id").Find(&res)
	for _, v := range res {
		ay.Db.Model(models.UserCard{}).Where("id = ?", v.Id).UpdateColumn("status", 2)
	}
}
