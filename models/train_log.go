/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package models

type TrainLog struct {
	BaseModel
	TrainId int64 `json:"train_id"`
	Uid     int64 `json:"uid"`
	Oid     int64 `json:"oid"`
}

func (TrainLog) TableName() string {
	return "d_train_log"
}
