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
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type LoginController struct {
}

type loginControllerLoginForm struct {
	Code     string `form:"code" binding:"required" label:"微信授权code"`
	Nickname string `form:"nickname" binding:"required" label:"姓名"`
	Avatar   string `form:"avatar" binding:"required" label:"头像"`
}

func (con LoginController) Login(c *gin.Context) {
	var data loginControllerLoginForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	isR, openid := con.GetOpenid(data.Code)

	if !isR {
		ay.Json{}.Msg(c, 400, openid, gin.H{})
	} else {

		user := models.UserModel{}.GetOpenid(openid)
		user.Avatar = data.Avatar
		user.Openid = openid
		user.NickName = data.Nickname
		ay.Db.Save(&user)

		ay.Json{}.Msg(c, 200, "success", gin.H{
			"token":    ay.AuthCode(strconv.Itoa(int(user.Id)), "ENCODE", "", 0),
			"avatar":   user.Avatar,
			"nickname": user.NickName,
		})
	}
}

func (con LoginController) GetOpenid(code string) (bool, string) {
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=" + ay.Yaml.GetString("wechat.appid") + "&secret=" + ay.Yaml.GetString("wechat.secret") + "&js_code=" + code + "&grant_type=authorization_code"

	response, err := http.Get(url)
	if err != nil || response.StatusCode != http.StatusOK {
		return false, string(rune(http.StatusServiceUnavailable))
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, err.Error()
	}
	defer response.Body.Close()

	type rj struct {
		Errcode    int    `json:"errcode"`
		Errmsg     string `json:"errmsg"`
		SessionKey string `json:"session_key"`
		Openid     string `json:"openid"`
	}

	var ms rj
	json.Unmarshal(body, &ms)
	log.Println(string(body))

	if ms.Errcode != 0 && ms.Openid != "" {
		return false, ms.Errmsg
	} else {
		return true, ms.Openid
	}
}
