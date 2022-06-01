/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package models

type ClassSpellingAskModel struct {
}

type ClassSpellingAsk struct {
	BaseModel
	SpellingId       int64  `json:"spelling_id"`
	UserName         string `json:"user_name"`
	UserAvatar       string `json:"user_avatar"`
	UserContent      string `json:"user_content"`
	EducationName    string `json:"education_name"`
	EducationContent string `json:"education_content"`
}

func (ClassSpellingAsk) TableName() string {
	return "d_class_spelling_ask"
}
