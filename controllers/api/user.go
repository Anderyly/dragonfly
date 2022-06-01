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
	"time"
)

type UserController struct {
}

// Info 用户信息
func (con UserController) Info(c *gin.Context) {

	isAuth, code, msg, user := Auth(GetToken(Token))

	if !isAuth {
		ay.Json{}.Msg(c, code, msg, gin.H{})
		return
	}

	isVip := 0

	if user.EffectiveAt.Unix() > time.Now().Unix() {
		isVip = 1
	}

	var vipLevel models.VipLevel

	vipLevelName := ay.Yaml.GetString("vip_default_level")
	vipNextLevelName := ""
	vipNextLevelNum := 0

	ay.Db.Order("num desc").
		Where("num <= ?", user.VipNum).
		First(&vipLevel)

	vipLevelName = vipLevel.Name
	var vipNextLevel models.VipLevel
	ay.Db.Order("num asc").
		Where("id > ?", vipLevel.Id).
		First(&vipNextLevel)
	if vipNextLevel.Id != 0 {
		vipNextLevelName = vipNextLevel.Name
		vipNextLevelNum = vipNextLevel.Num
	}

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"nickname":            user.NickName,
		"avatar":              user.Avatar,
		"amount":              user.Amount,
		"is_vip":              isVip,
		"vip_num":             user.VipNum,
		"vip_level":           vipLevelName,
		"effective_at":        user.EffectiveAt,
		"vip_next_level_name": vipNextLevelName,
		"vip_next_level_num":  vipNextLevelNum,
		"vip":                 vipLevel,
	})
}
