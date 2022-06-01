/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package models

type CardModel struct {
}

type Card struct {
	BaseModel
	Type          int     `json:"type"`
	Name          string  `json:"name"`
	Amount        float64 `json:"amount"`
	HourAmount    float64 `json:"hour_amount"`
	AllHour       int     `json:"all_hour"`
	EffectiveText string  `json:"effective_text"`
}

func (Card) TableName() string {
	return "d_card"
}
