/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package models

type TrainModel struct {
}

type Train struct {
	BaseModel
	StoreId     int64   `json:"store_id"`
	Room        string  `json:"room"`
	EducationId int64   `json:"education_id"`
	Name        string  `json:"name"`
	Address     string  `json:"address"`
	Mobile      string  `json:"mobile"`
	Introduce   string  `json:"introduce"`
	Attention   string  `json:"attention"`
	Arrange     string  `json:"arrange"`
	BackImage   string  `json:"back_image"`
	Image       string  `json:"image"`
	Video       string  `json:"video"`
	MinUser     int     `json:"min_user"`
	MaxUser     int     `json:"max_user"`
	StartDate   MyTime  `json:"start_date"`
	Amount      float64 `json:"amount"`
}

func (Train) TableName() string {
	return "d_train"
}
