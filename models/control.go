/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package models

type ControlModel struct {
}

type Control struct {
	BaseModel
	Type       int    `json:"type"`
	Cid        int64  `json:"cid"`
	Start      int64  `json:"start"`
	End        int64  `json:"end"`
	Status     int    `json:"status"`
	ControlId  string `json:"control_id"`
	ElectricSn string `json:"electric_sn"`
	Uid        int64  `json:"uid"`
}

func (Control) TableName() string {
	return "d_control"
}
