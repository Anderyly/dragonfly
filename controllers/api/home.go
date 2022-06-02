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
	"github.com/gin-gonic/gin"
)

type HomeController struct {
}

func (con HomeController) Banner(c *gin.Context) {

	var list []models.Banner
	ay.Db.Order("sort asc").Find(&list)

	for k, v := range list {
		list[k].Image = ay.Yaml.GetString("domain") + v.Image
		if v.Cover != "" {
			list[k].Cover = ay.Yaml.GetString("domain") + v.Cover
		}
	}

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"list": list,
	})
}
