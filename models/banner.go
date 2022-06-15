/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package models

type Banner struct {
	BaseModel
	Type  int    `json:"type"`
	Cover string `json:"cover"`
	Image string `json:"image"`
	Link  string `json:"link"`
	Sort  int    `json:"sort"`
}

func (Banner) TableName() string {
	return "d_banner"
}
