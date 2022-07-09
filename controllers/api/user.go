/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package api

import (
	"dragonfly/ay"
	"dragonfly/controllers"
	"dragonfly/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
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

	var orderCount int64
	ay.Db.Model(&models.Order{}).Where("uid = ?", user.Id).Count(&orderCount)

	var cardCount int64
	ay.Db.Model(&models.UserCard{}).Where("uid = ?", user.Id).Count(&cardCount)

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
		"order_num":           orderCount,
		"card_num":            cardCount,
	})
}

type orderUserControllerForm struct {
	Page int `form:"page" binding:"required" label:"页码"`
	Type int `form:"type" binding:"required" label:"类型"`
}

func (con UserController) Order(c *gin.Context) {
	var data orderUserControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	page := data.Page - 1

	isAuth, code, msg, user := Auth(GetToken(Token))

	if !isAuth {
		ay.Json{}.Msg(c, code, msg, gin.H{})
		return
	}

	var order []models.Order
	rs := ay.Db.Where("uid = ?", user.Id).
		Order("created_at desc").
		Offset(page * Limit).
		Limit(Limit)

	today := time.Now().Format("20060102")
	hour := time.Now().Format("15")

	switch data.Type {
	case 2:
		rs.Where("status = 2")
	case 3:
		rs.Where("content like ? AND content like ? AND status = 1", "%"+today+"%", "%"+hour+"%")
	case 4:
		rs.Where("content not like ? AND content not like ? AND status = 1", "%"+today+"%", "%"+hour+"%")
	default:
		rs.Where("status = 1 OR status = 2")
	}

	rs.Debug().Find(&order)

	var list []gin.H
	for _, v := range order {
		list = append(list, gin.H{
			"id":         v.Id,
			"type":       v.Type,
			"cid":        v.Cid,
			"status":     v.Status,
			"op":         v.Op,
			"created_at": v.CreatedAt.Unix(),
			"amount":     v.Amount,
		})
	}

	if list == nil {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"list": []string{},
		})
	} else {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"list": list,
		})
	}
}

func (con UserController) GetQr(c *gin.Context) {

	isAuth, code, msg, user := Auth(GetToken(Token))

	if !isAuth {
		ay.Json{}.Msg(c, code, msg, gin.H{})
		return
	}

	t := time.Now().Unix()

	var res models.Control
	ay.Db.Where("uid = ? AND type = 1 AND ((? - start <= 1800) or end > ?) AND end > 0", user.Id, t, t).Order("start asc").First(&res)

	if res.Id == 0 {
		ay.Json{}.Msg(c, 400, "未到时间", gin.H{})
	} else {

		var room models.Room
		ay.Db.Select("id,name").Where("control_id = ?", res.ControlId).First(&room)

		h := time.Unix(res.Start, 0).Format("15")

		var order models.Order
		ay.Db.Debug().Where("uid = ? AND cid = ? AND content like ? AND content like ?", res.Uid, room.Id, "%"+h+"%", "%"+time.Unix(res.Start, 0).Format("20060102")+"%").First(&order)

		type cv struct {
			Num  int      `json:"num"`
			Time []string `json:"time"`
			Ymd  int      `json:"ymd"`
		}
		var cc cv
		json.Unmarshal([]byte(order.Content), &cc)

		id := res.ControlId + "," + ay.Yaml.GetString("public_control_id")

		sT := time.Unix(res.Start, 0)
		eT := time.Unix(res.End, 0)

		controllers.ControlServer{}.EditUser(user.ControlUserId, sT.Format("2006-01-02 15:04:05"), eT.Format("2006-01-02 15:04:05"))
		controllers.ControlServer{}.BindUser(user.ControlUserId, id, "")
		_, text := controllers.ControlServer{}.GetQr(user.ControlUserId)
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"text":           text,
			"room_name":      room.Name,
			"subscribe_time": cc.Time,
			"ymd":            cc.Ymd,
		})

	}
}

func (con UserController) Open(c *gin.Context) {
	isAuth, code, msg, user := Auth(GetToken(Token))

	if !isAuth {
		ay.Json{}.Msg(c, code, msg, gin.H{})
		return
	}

	t := time.Now().Unix()

	var res models.Control
	ay.Db.Where("uid = ? AND type = 1 AND ((? - start <= 15*60) or end > ?) AND end > 0", user.Id, t, t).Order("start asc").First(&res)

	if res.Id == 0 {
		ay.Json{}.Msg(c, 400, "未到时间", gin.H{})
	} else {
		id := controllers.ControlServer{}.GetId(ay.Yaml.GetString("public_control_sn"))
		log.Println(id)
		res := controllers.ControlServer{}.Set(strconv.Itoa(id), "10")
		log.Println(res)
		ay.Json{}.Msg(c, 200, "已开门", gin.H{})

	}
}
