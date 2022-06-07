/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package models

type UserChannel struct {
	BaseModel
	Type int   `json:"type"`
	Uid  int64 `json:"uid"`
	Cid  int64 `json:"cid"`
	Oid  int64 `json:"oid"`
}

func (UserChannel) TableName() string {
	return "d_user_channel"
}
