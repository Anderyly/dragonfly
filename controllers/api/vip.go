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

type VipController struct {
}

func (con VipController) Amount(c *gin.Context) {

	var vip []models.Vip

	ay.Db.Order("sort asc").Find(&vip)

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"list": vip,
	})
}

func (con VipController) Pay(c *gin.Context) {

}
