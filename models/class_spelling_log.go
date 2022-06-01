/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package models

type ClassSpellingLogModel struct {
}

type ClassSpellingLog struct {
	BaseModel
	SpellingId int64 `json:"spelling_id"`
	Uid        int64 `json:"uid"`
}

func (ClassSpellingLog) TableName() string {
	return "d_class_spelling_log"
}
