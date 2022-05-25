package api

import (
	"dragonfly/ay"
	"dragonfly/models"
	"fmt"
	"github.com/gin-gonic/gin"
	qrcode "github.com/skip2/go-qrcode"
	"time"
)

var (
	Token string
	Limit = 10
)

type CommonController struct {
}

func (con CommonController) A() {

}

func GetRequestIP(c *gin.Context) string {
	reqIP := c.ClientIP()
	if reqIP == "::1" {
		reqIP = "127.0.0.1"
	}
	return reqIP
}

func GetToken(token string) string {
	uid := ay.AuthCode(token, "DECODE", "", 0)
	return uid
}

func Auth(uid string) (bool, int, string, models.User) {
	var user models.User
	ay.Db.First(&user, "id = ?", uid)

	if user.Id == 0 {
		return false, 401, "token错误", user
	}
	return true, 200, "success", user

}

func (con CommonController) MakeQrCode(text string) string {
	name := ay.MD5(fmt.Sprintf("%s%s", text, time.Now().String())) + ".png"
	fileDir := fmt.Sprintf("static/qrcode/%d-%d/", time.Now().Year(), time.Now().Month())

	err := ay.CreateMutiDir(fileDir)
	if err != nil {
		return ""
	}

	err = qrcode.WriteFile(text, qrcode.Medium, 152, fileDir+name)
	if err != nil {
		return ""
	}
	return fileDir + name

}
