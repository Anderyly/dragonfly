/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package task

import (
	"dragonfly/ay"
	"dragonfly/controllers"
	"dragonfly/models"
	"log"
	"strconv"
	"time"
)

var sjd = [24]string{
	"01:00",
	"02:00",
	"03:00",
	"04:00",
	"05:00",
	"06:00",
	"07:00",
	"08:00",
	"09:00",
	"10:00",
	"11:00",
	"12:00",
	"13:00",
	"14:00",
	"15:00",
	"16:00",
	"17:00",
	"18:00",
	"19:00",
	"20:00",
	"21:00",
	"22:00",
	"23:00",
	"00:00",
}

func Control() {
	// TODO
	/**
	修改修改：
		1. 用户预约将绑定门禁(教室门禁以及公共区域门禁)
		2. 监听用户到指定预约时间开启电
		2. 监听用户到期时间关门关电（教室门禁5分钟 公共区域延迟1个半小时）
	*/
	log.Println("监听开门等操作")

	var subscribe models.UserRoomSubscribe
	ay.Db.Where("status = 0 and ymd = ?", time.Now().Format("20060102")).Order("hour_id asc").First(&subscribe)

	if subscribe.Id == 0 {
		return
	}

	// 预约时间戳
	ymd := strconv.Itoa(subscribe.Ymd)
	vt := ymd[0:4] + "-" + ymd[4:6] + "-" + ymd[6:8] + " " + sjd[subscribe.HourId-1] + ":00"
	log.Println(vt)
	stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", vt, time.Local)
	subT := stamp.Unix()
	log.Println(subT)

	var ids []int64
	ids = append(ids, subscribe.Id)

	// 关门时间 预先一小时
	var closeDoor int64 = 3600

	var closeElectric int64 = 3600

	// 查找相邻的时间
	var res []models.UserRoomSubscribe
	ay.Db.Select("id,hour_id,uid").
		Where("status = 0 and ymd = ? and id != ?", time.Now().Format("20060102"), subscribe.Id).
		Order("hour_id asc").
		Find(&res)
	hour := subscribe.HourId
	for _, v := range res {
		if v.HourId != hour+1 {
			break
		} else {
			// 当前用户使用加一小时关门时间
			if v.Uid == subscribe.Uid {
				closeDoor += 3600
				ids = append(ids, v.Id)
			}
			// 一直有使用将一直加关闭电源
			closeElectric += 3600
			hour = v.HourId
		}
	}
	log.Println(ids)

	closeDoor += subT
	closeElectric += subT

	var room models.Room
	ay.Db.First(&room, subscribe.RoomId)
	// 教室 开门
	ay.Db.Create(&models.Control{
		Cid: subscribe.Id,

		Type:      1,
		Start:     subT,
		End:       closeElectric,
		Status:    0,
		ControlId: room.ControlId,
		//ElectricSn: room.Electric,
		Uid: subscribe.Uid,
	})

	// 教室关门
	ay.Db.Create(&models.Control{
		Cid: subscribe.Id,

		Type: 2,
		//Start: subT,
		End:       closeElectric,
		Status:    0,
		ControlId: room.ControlId,
		//ElectricSn: room.Electric,
		Uid: subscribe.Uid,
	})

	// 公共区域 开门
	ay.Db.Create(&models.Control{
		Cid: subscribe.Id,

		Type:  1,
		Start: subT,
		//End:       closeElectric,
		Status:    0,
		ControlId: ay.Yaml.GetString("public_control_id"),
		//ElectricSn: room.Electric,
		Uid: subscribe.Uid,
	})

	// 公共区域关门
	ay.Db.Create(&models.Control{
		Cid:  subscribe.Id,
		Type: 2,
		//Start: subT,
		End:       subT + 15*60,
		Status:    0,
		ControlId: ay.Yaml.GetString("public_control_id"),
		//ElectricSn: room.Electric,
		Uid: subscribe.Uid,
	})

	// 教室 开电
	ay.Db.Create(&models.Control{
		Cid: subscribe.Id,

		Type:  3,
		Start: subT,
		//End:        closeElectric,
		Status: 0,
		//ControlId: room.ControlId,
		ElectricSn: room.Electric,
		Uid:        subscribe.Uid,
	})

	// 教室 关电
	ay.Db.Create(&models.Control{
		Cid: subscribe.Id,

		Type: 4,
		//Start: subT,
		End:    closeElectric,
		Status: 0,
		//ControlId: room.ControlId,
		ElectricSn: room.Electric,
		Uid:        subscribe.Uid,
	})

	for _, v := range ids {
		ay.Db.Model(models.UserRoomSubscribe{}).Where("id = ?", v).UpdateColumn("status", 1)
	}

	//subscribe.Status = 1
	//ay.Db.Save(&subscribe)
}

func ControlDel() {
	var res []models.Control
	ay.Db.Select("id,cid").Group("cid").Find(&res)
	for _, v := range res {
		var rs models.UserRoomSubscribe
		ay.Db.First(&rs, v.Cid)
		if rs.Id == 0 {
			ay.Db.Where("cid = ?", v.Cid).Delete(&models.Control{})
		}
	}
}

// ControlOpenElectric 开电
func ControlOpenElectric() {
	t := time.Now().Unix()
	var res models.Control
	ay.Db.Where("type = 3 AND status = 0 AND (start = ? or (start - ? < 300 AND start - ? > 0))", t, t, t).Select("id,uid,electric_sn").First(&res)

	if res.Id == 0 {
		return
	}
	controllers.ElectricServer{}.Set(1, res.ElectricSn)

	ay.Db.Model(models.Control{}).Where("id = ?", res.Id).UpdateColumn("status", 1)

}

// ControlCloseElectric 关电
func ControlCloseElectric() {
	t := time.Now().Unix()
	var res models.Control
	ay.Db.Where("type = 4 AND status = 0 AND ? - end > 300", t).Select("id,uid,electric_sn").First(&res)

	if res.Id == 0 {
		return
	}
	controllers.ElectricServer{}.Set(0, res.ElectricSn)

	ay.Db.Model(models.Control{}).Where("id = ?", res.Id).UpdateColumn("status", 1)
}

// ControlCloseDoor 关门
func ControlCloseDoor() {

	t := time.Now().Unix()
	var res models.Control
	ay.Db.Where("type = 2 AND status = 0 AND ? - end > 300", t).Select("id,uid,control_id").First(&res)

	if res.Id == 0 {
		return
	}
	controllers.ControlServer{}.Set(res.ControlId, "1")

	ay.Db.Model(models.Control{}).Where("id = ?", res.Id).UpdateColumn("status", 1)
}
