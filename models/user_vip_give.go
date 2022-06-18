/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package models

type UserVipGiveModel struct {
}

type UserVipGive struct {
	BaseModel
	Uid int64 `json:"uid"`
	Ymd int   `json:"ymd"`
}

func (UserVipGive) TableName() string {
	return "d_user_vip_give"
}
