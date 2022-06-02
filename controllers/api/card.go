/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package api

import (
	"dragonfly/ay"
	"dragonfly/models"
	"github.com/gin-gonic/gin"
)

type CardController struct {
}

type listCardControllerForm struct {
	Page int `form:"page" binding:"required" label:"页码"`
}

func (con CardController) List(c *gin.Context) {
	var data listCardControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	page := data.Page - 1

	var coupon []models.Card
	ay.Db.Offset(page * Limit).
		Limit(Limit).Debug().
		Order("created_at asc").
		Find(&coupon)

	var list []gin.H
	for _, v := range coupon {
		list = append(list, gin.H{
			"id":             v.Id,
			"type":           v.Type,
			"name":           v.Name,
			"amount":         v.Amount,
			"hour_amount":    v.HourAmount,
			"all_hour":       v.AllHour,
			"effective_text": v.EffectiveText,
		})
	}

	if list == nil {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"list": []string{},
		})
	} else {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"list": list,
		})
	}
}

type userCardControllerForm struct {
	Page int `form:"page" binding:"required" label:"页码"`
}

func (con CardController) User(c *gin.Context) {
	var data userCardControllerForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	isAuth, code, msg, requestUser := Auth(GetToken(Token))

	if !isAuth {
		ay.Json{}.Msg(c, code, msg, gin.H{})
		return
	}

	page := data.Page - 1

	type r struct {
		Id          int64         `json:"id"`
		Type        int           `json:"type"`
		Name        string        `json:"name"`
		Amount      float64       `json:"amount"`
		HourAmount  float64       `json:"hour_amount"`
		AllHour     int           `json:"all_hour"`
		EffectiveAt models.MyTime `json:"effective_at"`
		UseHour     int           `json:"use_hour"`
		Status      int           `json:"status"`
	}
	var res []r
	ay.Db.Table("d_user_card").
		Select("d_card.type,d_card.name,d_card.amount,d_card.hour_amount,d_card.all_hour,d_user_card.use_hour,d_user_card.status,d_user_card.effective_at").
		Joins("left join d_card on d_user_card.card_id=d_card.id").
		Offset(page*Limit).
		Limit(Limit).Debug().
		Where("d_user_card.uid = ?", requestUser.Id).
		Order("d_user_card.status asc").
		Find(&res)

	var list []gin.H
	for _, v := range res {
		hour := v.AllHour - v.UseHour
		list = append(list, gin.H{
			"id":           v.Id,
			"type":         v.Type,
			"name":         v.Name,
			"amount":       v.Amount,
			"hour_amount":  v.HourAmount,
			"all_hour":     hour,
			"effective_at": v.EffectiveAt.Unix(),
			"status":       v.Status,
		})
	}

	if list == nil {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"list": []string{},
		})
	} else {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"list": list,
		})
	}
}

// Pay 会员卡办理
func (con CardController) Pay(c *gin.Context) {

}
