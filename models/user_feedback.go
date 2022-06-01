/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package models

type FeedbackModel struct {
}

type Feedback struct {
	BaseModel
	Uid     int64  `json:"uid"`
	Name    string `json:"name"`
	Mobile  string `json:"mobile"`
	Content string `json:"content"`
	Status  int    `json:"status"`
}

func (Feedback) TableName() string {
	return "d_user_feedback"
}
