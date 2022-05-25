/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package ay

import "gorm.io/gorm"

type Base struct {
}

var (
	Db *gorm.DB
)

func init() {

}
