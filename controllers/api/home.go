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
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type HomeController struct {
}

func (con HomeController) Main(c *gin.Context) {
	var config models.Config
	ay.Db.First(&config, 1)

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"info": config,
	})
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

// Chart 排行榜
func (con HomeController) Chart(c *gin.Context) {

	_, _, _, requestUser := Auth(GetToken(Token))

	type r struct {
		Avatar   string  `json:"avatar"`
		Nickname string  `json:"nickname"`
		Uid      int64   `json:"uid"`
		Count    int64   `json:"count"`
		Time     float64 `json:"time"`
	}

	where := ""

	if c.Query("type") == "1" {
		where = "t1.ymd = " + time.Now().Format("20060102")
	} else if c.Query("type") == "2" {
		where = "t1.ymd > " + fmt.Sprintf("%s", time.Now().Add(-24*time.Hour*7).Format("20060102")) + " AND t1.ymd <= " + time.Now().Format("20060102")
	} else {
		where = "t1.ymd like '" + time.Now().Format("200601") + "%'"
	}

	var rj []r

	sql := `SELECT DISTINCT (
SELECT uid
FROM d_user_room_subscribe t2
WHERE t1.uid = t2.uid
GROUP BY uid
ORDER BY count(*) desc
LIMIT 1) as uid,(
SELECT count(hour_id) 
FROM d_user_room_subscribe t4
where t4.uid=t1.uid
) as time, t3.avatar, t3.nickname
FROM d_user_room_subscribe t1 join d_user as t3 on t1.uid=t3.id WHERE ` + where + " limit 10"

	ay.Db.Raw(sql).Scan(&rj)
	var list []gin.H
	for _, v := range rj {
		var count int64
		ay.Db.Model(&models.Order{}).Where("type = 1 and uid = ? and status = 1", v.Uid).Count(&count)
		list = append(list, gin.H{
			"nickname": v.Nickname,
			"avatar":   v.Avatar,
			"count":    count,
			"time":     v.Time,
		})
	}

	var count int64
	ay.Db.Model(&models.Order{}).Where("type = 1 and uid = ? and status = 1", requestUser.Id).Count(&count)

	var t int64
	ay.Db.Model(&models.UserRoomSubscribe{}).Where("uid = ?", requestUser.Id).Count(&t)

	var average int64
	if t == 0 || count == 0 {
		average = 0
	} else {
		average = t / count
	}
	if list == nil {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"list": []string{},
			"user": gin.H{
				"count":   count,
				"time":    t,
				"average": average,
			},
		})
	} else {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"list": list,
			"user": gin.H{
				"count":   count,
				"time":    t,
				"average": average,
			},
		})
	}

}
