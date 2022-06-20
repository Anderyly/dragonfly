/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package models

type BranchStoreModel struct {
}

type BranchStore struct {
	BaseModel
	Name          string `json:"name"`
	Address       string `json:"address"`
	Mobile        string `json:"mobile"`
	Introduce     string `json:"introduce"`
	Attention     string `json:"attention"`
	Image         string `json:"image"`
	Video         string `json:"video"`
	Lat           string `json:"lat"`
	Lng           string `json:"lng"`
	Label         string `json:"label"`
	ControlAreaId string `json:"control_area_id"`
}

func (BranchStore) TableName() string {
	return "d_branch_store"
}
