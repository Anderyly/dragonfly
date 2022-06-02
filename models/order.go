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
	Type        int     `json:"type"`
	Cid         int64   `json:"cid"`
	Uid         int64   `json:"uid"`
	OutTradeNo  string  `json:"out_trade_no"`
	TradeNo     string  `json:"trade_no"`
	Amount      float64 `json:"amount"`
	OldAmount   float64 `json:"old_amount"`
	CouponId    int64   `json:"coupon_id"`
	CardId      int64   `json:"card_id"`
	VipDiscount float64 `json:"vip_discount"`
	PayAt       MyTime  `json:"pay_at"`
	PayType     int     `json:"pay_type"`
	Status      int     `json:"status"`
	Op          int     `json:"op"`
	Content     string  `json:"content"`
	Des         string  `json:"des"`
	Ip          string  `json:"ip"`
}

func (Order) TableName() string {
	return "d_order"
}
