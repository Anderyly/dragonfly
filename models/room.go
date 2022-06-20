/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package models

type RoomModel struct {
}

type Room struct {
	BaseModel
	StoreId   int64   `json:"store_id"`
	Type      int     `json:"type"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Mobile    string  `json:"mobile"`
	Introduce string  `json:"introduce"`
	Attention string  `json:"attention"`
	Guide     string  `json:"guide"`
	Cancel    string  `json:"cancel"`
	Image     string  `json:"image"`
	Video     string  `json:"video"`
	Size      float64 `json:"size"`
	Electric  string  `json:"electric"`
	ControlSn string  `json:"control_sn"`
	ControlId string  `json:"control_id"`
	DeptId    string  `json:"dept_id"`
}

func (Room) TableName() string {
	return "d_room"
}
