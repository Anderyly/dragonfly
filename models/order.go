/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package models

type OrderModel struct {
}

type Order struct {
	BaseModel
	Oid        string  `gorm:"column:oid" json:"oid"`
	OutTradeNo string  `gorm:"column:out_trade_no" json:"out_trade_no"`
	Type       int     `gorm:"column:type" json:"type"`
	Uid        int64   `gorm:"column:uid" json:"uid"`
	ToUid      int64   `gorm:"column:to_uid" json:"to_uid"`
	UserName   string  `gorm:"column:username" json:"username"`
	Gender     int     `gorm:"column:gender" json:"gender"`
	Y          int     `gorm:"column:y" json:"y"`
	M          int     `gorm:"column:m" json:"m"`
	D          int     `gorm:"column:d" json:"d"`
	H          int     `gorm:"column:h" json:"h"`
	I          int     `gorm:"column:i" json:"i"`
	S          int     `gorm:"column:s" json:"s"`
	Des        string  `gorm:"column:des" json:"des"`
	Remark     string  `gorm:"column:remark" json:"remark"`
	Amount     float64 `gorm:"column:amount" json:"amount"`
	OldAmount  float64 `gorm:"column:old_amount" json:"old_amount"`
	Coupon     int64   `gorm:"column:coupon" json:"coupon"`
	Line       string  `gorm:"column:line" json:"line"`
	Ip         string  `gorm:"column:ip" json:"ip"`
	Appid      int64   `gorm:"column:appid" json:"appid"`
	PayType    int     `gorm:"column:pay_type" json:"pay_type"`
	TradeNo    string  `gorm:"column:trade_no" json:"trade_no"`
	AreaId     int     `gorm:"column:area_id" json:"area_id"`
	Status     int     `gorm:"column:status" json:"status"`
	Op         int     `gorm:"column:op" json:"op"`
	ReturnUrl  string  `gorm:"column:return_url" json:"-"`
	Json       string  `gorm:"column:json" json:"json"`
	PayTime    string  `gorm:"column:pay_time" json:"pay_time"`
	Birthday   string  `gorm:"column:birthday" json:"birthday"`
	Content    string  `gorm:"column:content" json:"content"`
	Discount   float64 `gorm:"column:discount" json:"discount"`
}

func (Order) TableName() string {
	return "d_order"
}
