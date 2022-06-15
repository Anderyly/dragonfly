package admin

import (
	"dragonfly/ay"
	"dragonfly/controllers/api"
	"dragonfly/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"path"
	"strings"
	"time"
)

type CommonController struct {
}

func Auth() bool {
	token := ay.AuthCode(api.Token, "DECODE", "", 0)
	if token == "" {
		log.Println("没有token")
		return false
	}
	var user models.Admin
	ay.Db.First(&user, "account = ?", token)
	if user.Id == 0 {
		return false
	} else {
		return true
	}
}

func Upload(c *gin.Context, address string) (int, string, string) {
	file, err := c.FormFile("file")
	if err != nil {
		return 400, err.Error(), ""

	}

	fileExt := strings.ToLower(path.Ext(file.Filename))
	log.Println(fileExt)
	if fileExt != ".png" && fileExt != ".jpg" && fileExt != ".gif" && fileExt != ".jpeg" && fileExt != ".doc" && fileExt != ".docx" && fileExt != ".pdf" && fileExt != ".mp4" {
		return 400, "上传失败!只允许png,jpg,gif,jpeg,doc,pdf,docx,mp4文件", ""
	}
	fileName := ay.MD5(fmt.Sprintf("%s%s", file.Filename, time.Now().String()))
	fildDir := fmt.Sprintf("static/upload/admin/"+address+"/%d-%d/", time.Now().Year(), time.Now().Month())

	err = ay.CreateMutiDir(fildDir)
	if err != nil {
		return 400, err.Error(), ""

	}

	filepath := fmt.Sprintf("%s%s%s", fildDir, fileName, fileExt)
	c.SaveUploadedFile(file, filepath)

	return 200, "/" + filepath, fileName

}
