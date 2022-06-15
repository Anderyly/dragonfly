/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package models

type AdminModel struct {
}

type Admin struct {
	Id       int64
	Account  string
	Password string
}

func (Admin) TableName() string {
	return "d_admin"
}
