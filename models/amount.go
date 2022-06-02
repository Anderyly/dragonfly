/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package models

type Amount struct {
	BaseModel
	Amount     float64 `json:"amount"`
	GiveAmount float64 `json:"give_amount"`
	Sort       int     `json:"sort"`
}

func (Amount) TableName() string {
	return "d_amount"
}
