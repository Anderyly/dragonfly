/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package api

import (
	"dragonfly/ay"
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

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"nickname":     user.NickName,
		"avatar":       user.Avatar,
		"amount":       user.Amount,
		"is_vip":       isVip,
		"vip_num":      user.VipNum,
		"effective_at": user.EffectiveAt,
	})
}
