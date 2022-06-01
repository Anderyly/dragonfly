/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package models

type EducationModel struct {
}

type Education struct {
	BaseModel
	StoreId   int64   `json:"store_id"`
	Name      string  `json:"name"`
	Avatar    string  `json:"avatar"`
	Label     string  `json:"label"`
	Introduce string  `json:"introduce"`
	Amount    float64 `json:"amount"`
	Sign      string  `json:"sign"`
	Wechat    string  `json:"wechat"`
	Image     string  `json:"image"`
	Video     string  `json:"video"`
}

func (Education) TableName() string {
	return "d_education"
}
