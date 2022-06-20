/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package controllers

import (
	"dragonfly/ay"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ControlServer struct {
}

var (
	appid  = "f910c17100764058bac294019d5f16d0"
	secret = "c64b2398c73e42fbabbac9302ca4ba592a994f4457a344a98af5570daf89a4b0df9eeef1547745a7bd18dc532b3e606d96e33e1f3551483c93c132ac661f0c83"
)

type returnRes struct {
	Ret     int    `json:"ret"`
	ErrMsg  string `json:"errMsg"`
	Data    bool   `json:"data"`
	Success bool   `json:"success"`
}

// GetCard 生成卡号
func (con ControlServer) GetCard() (bool, string) {
	sign, t := con.getSign()

	paths := "/open/persion/getRandomCard?appid=" + appid + "&timestamp=" + t + "&sign=" + sign
	res := con.Http(paths, "")
	type es struct {
		Ret     int         `json:"ret"`
		ErrMsg  interface{} `json:"errMsg"`
		Data    string      `json:"data"`
		Success bool        `json:"success"`
	}
	var rj es
	json.Unmarshal([]byte(res), &rj)
	if rj.Success {
		return true, rj.Data
	} else {
		return false, ""
	}
}

// BindUser 用户绑定门禁
func (con ControlServer) BindUser(id, add, del string) bool {
	sign, t := con.getSign()

	paths := "/open/persion/bindPersionDoor?appid=" + appid + "&timestamp=" + t + "&sign=" + sign + "&persionId=" + id

	if del != "" {
		paths += "&delDoorIdList=" + del
	}
	if add != "" {
		paths += "&addDoorIdList=" + add
	}
	res := con.Http(paths, "")
	var rj returnRes
	json.Unmarshal([]byte(res), &rj)
	if rj.Success {
		return true
	} else {
		return false
	}
}

// GetQr 获取二维码
func (con ControlServer) GetQr(id string) (bool, string) {
	sign, t := con.getSign()

	paths := "/open/persion/getQrCode?appid=" + appid + "&timestamp=" + t + "&sign=" + sign + "&persionId=" + id
	res := con.Http(paths, "")
	type es struct {
		Ret    int         `json:"ret"`
		ErrMsg interface{} `json:"errMsg"`
		Data   struct {
			CreateTime interface{} `json:"createTime"`
			ModifyTime interface{} `json:"modifyTime"`
			Card       int         `json:"card"`
			Time       int         `json:"time"`
			Sign       string      `json:"sign"`
		} `json:"data"`
		Success bool `json:"success"`
	}
	var rj es
	json.Unmarshal([]byte(res), &rj)
	if rj.Success {
		return true, rj.Data.Sign
	} else {
		return false, ""
	}
}

// AddUser 新增用户
func (con ControlServer) AddUser(name, mobile, start, end, dept, status string) (bool, int) {
	_, card := con.GetCard()
	sign, t := con.getSign()
	param := `{
"canAttence": true,
"canPass": true,
"card": "` + card + `",
"deptIdList": [` + dept + `],
"effectStartTime": "` + start + `",
"effectEndTime": "` + end + `",
"mobile": "` + mobile + `",
"name": "` + name + `",
"openType": 4,
"persionType": 1,
"faceImgUrl":"http://tutk-bull-test.oss-cn-hangzhou.aliyuncs.com/20190904_44f876a32aab48e4a8f79359cfb39516.png",
"status": ` + status + `
}`
	paths := "/open/persion/insert?appid=" + appid + "&timestamp=" + t + "&sign=" + sign
	res := con.Http(paths, param)

	type rr struct {
		Ret    int         `json:"ret"`
		ErrMsg interface{} `json:"errMsg"`
		Data   struct {
			CreateTime      int64         `json:"createTime"`
			ModifyTime      interface{}   `json:"modifyTime"`
			Id              int           `json:"id"`
			Name            string        `json:"name"`
			Mobile          string        `json:"mobile"`
			Number          interface{}   `json:"number"`
			Card            int           `json:"card"`
			EffectStartTime string        `json:"effectStartTime"`
			EffectEndTime   string        `json:"effectEndTime"`
			Status          int           `json:"status"`
			StatusStr       interface{}   `json:"statusStr"`
			OrgId           int           `json:"orgId"`
			DeptId          interface{}   `json:"deptId"`
			OpenType        int           `json:"openType"`
			OpenTypeStr     interface{}   `json:"openTypeStr"`
			DeptIdList      []int         `json:"deptIdList"`
			DoorList        []interface{} `json:"doorList"`
			DeptList        []interface{} `json:"deptList"`
			PersionType     int           `json:"persionType"`
			FaceImgUrl      string        `json:"faceImgUrl"`
			PersionTypeStr  interface{}   `json:"persionTypeStr"`
			AddDoorIdList   []interface{} `json:"addDoorIdList"`
			ParentNamePath  interface{}   `json:"parentNamePath"`
			AuthDoorName    interface{}   `json:"authDoorName"`
		} `json:"data"`
		Success bool `json:"success"`
	}
	var rj rr
	json.Unmarshal([]byte(res), &rj)
	if rj.Success {
		return true, rj.Data.Id
	} else {
		return false, 0
	}
}

// EditUser 编辑用户
func (con ControlServer) EditUser(id, start, end string) bool {
	_, card := con.GetCard()
	sign, t := con.getSign()
	param := `{
"id":` + id + `,
"canAttence": true,
"card": "` + card + `",
"canPass": true,
"effectStartTime": "` + start + `",
"effectEndTime": "` + end + `"
}`
	paths := "/open/persion/editById?appid=" + appid + "&timestamp=" + t + "&sign=" + sign
	res := con.Http(paths, param)

	var rj returnRes
	json.Unmarshal([]byte(res), &rj)
	if rj.Success {
		return true
	} else {
		return false
	}
}

// AddDept 新增部门
func (con ControlServer) AddDept(name, pid string) (bool, int) {
	sign, t := con.getSign()
	param := `{
"name": "` + name + `",
"parentId": ` + pid + `,
"remark": "` + name + `"
}`
	paths := "/open/dept/insert?appid=" + appid + "&timestamp=" + t + "&sign=" + sign
	res := con.Http(paths, param)
	var rj returnRes
	json.Unmarshal([]byte(res), &rj)
	if rj.Success {
		id := con.GetDeptId(name)
		return true, id
	} else {
		return false, 0
	}
}

// EditDept 修改部门
func (con ControlServer) EditDept(name, id string) bool {
	sign, t := con.getSign()
	param := `{
"id": "` + id + `",
"name": "` + name + `",
"remark": "` + name + `"
}`
	paths := "/open/dept/editById?appid=" + appid + "&timestamp=" + t + "&sign=" + sign
	res := con.Http(paths, param)

	var rj returnRes
	json.Unmarshal([]byte(res), &rj)
	if rj.Success {
		return true
	} else {
		return false
	}
}

// GetDeptId 获取部门id
func (con ControlServer) GetDeptId(name string) int {
	sign, t := con.getSign()
	param := `{
"name": "` + name + `"
}`
	paths := "/open/dept/getList?appid=" + appid + "&timestamp=" + t + "&sign=" + sign
	res := con.Http(paths, param)

	type r struct {
		Ret    int    `json:"ret"`
		ErrMsg string `json:"errMsg"`
		Data   []struct {
			CreateTime    int64         `json:"createTime"`
			ModifyTime    string        `json:"modifyTime"`
			Id            int           `json:"id"`
			Name          string        `json:"name"`
			Disabled      bool          `json:"disabled"`
			OrgId         int           `json:"orgId"`
			Remark        string        `json:"remark"`
			ParentId      int           `json:"parentId"`
			ChildList     []interface{} `json:"childList"`
			PersionIdList string        `json:"persionIdList"`
		} `json:"data"`
		Ss bool `json:"ss"`
	}
	var rj r
	json.Unmarshal([]byte(res), &rj)
	for _, v := range rj.Data {
		if v.Name == name {
			return v.Id
		}
	}
	return 0
}

// DelDept 删除部门
func (con ControlServer) DelDept(id string) bool {
	sign, t := con.getSign()
	paths := "/open/dept/deleteById?appid=" + appid + "&timestamp=" + t + "&sign=" + sign + "&id=" + id
	res := con.Http(paths, "")

	var rj returnRes
	json.Unmarshal([]byte(res), &rj)
	if rj.Success {
		return true
	} else {
		return false
	}
}

// BindDevice 绑定区域
func (con ControlServer) BindDevice(id string, areaId string) bool {
	sign, t := con.getSign()
	paths := "/open/area/bindAreaDevice?appid=" + appid + "&timestamp=" + t + "&sign=" + sign + "&areaId=" + areaId + "&deviceIdList=" + id
	res := con.Http(paths, "")
	var rj returnRes
	json.Unmarshal([]byte(res), &rj)
	if rj.Success {
		return true
	} else {
		return false
	}
}

// AddArea 新增区域
func (con ControlServer) AddArea(name, pid string) (bool, int) {
	sign, t := con.getSign()
	param := `{
"name": "` + name + `",
"parentId": ` + pid + `,
"remark": "` + name + `"
}`
	paths := "/open/area/insert?appid=" + appid + "&timestamp=" + t + "&sign=" + sign
	res := con.Http(paths, param)

	var rj returnRes
	json.Unmarshal([]byte(res), &rj)
	if rj.Success {
		id := con.GetAreaId(name)
		return true, id
	} else {
		return false, 0
	}
}

// EditArea 修改区域
func (con ControlServer) EditArea(name, id string) bool {
	sign, t := con.getSign()
	param := `{
"id": "` + id + `",
"name": "` + name + `",
"remark": "` + name + `"
}`
	paths := "/open/area/editById?appid=" + appid + "&timestamp=" + t + "&sign=" + sign
	res := con.Http(paths, param)

	var rj returnRes
	json.Unmarshal([]byte(res), &rj)
	if rj.Success {
		return true
	} else {
		return false
	}
}

// DelArea 删除区域
func (con ControlServer) DelArea(id string) bool {
	sign, t := con.getSign()
	paths := "/open/area/deleteById?appid=" + appid + "&timestamp=" + t + "&sign=" + sign + "&id=" + id
	res := con.Http(paths, "")

	var rj returnRes
	json.Unmarshal([]byte(res), &rj)
	if rj.Success {
		return true
	} else {
		return false
	}
}

// GetAreaId 获取区域id
func (con ControlServer) GetAreaId(name string) int {
	sign, t := con.getSign()
	param := `{
"name": "` + name + `"
}`
	paths := "/open/area/getList?appid=" + appid + "&timestamp=" + t + "&sign=" + sign
	res := con.Http(paths, param)

	type r struct {
		Ret    int         `json:"ret"`
		ErrMsg interface{} `json:"errMsg"`
		Data   []struct {
			CreateTime   int64         `json:"createTime"`
			ModifyTime   interface{}   `json:"modifyTime"`
			Id           int           `json:"id"`
			Name         string        `json:"name"`
			OrgId        int           `json:"orgId"`
			Remark       string        `json:"remark"`
			Disabled     bool          `json:"disabled"`
			ParentId     int           `json:"parentId"`
			ChildList    []interface{} `json:"childList"`
			DeviceIdlist interface{}   `json:"deviceIdlist"`
		} `json:"data"`
		Success bool `json:"success"`
	}

	var rj r
	json.Unmarshal([]byte(res), &rj)
	for _, v := range rj.Data {
		if v.Name == name {
			return v.Id
		}
	}
	return 0
}

// Set 设置设备
func (con ControlServer) Set(sn, timestamp string) string {
	id := con.GetId(sn)
	//id := 6204
	sign, t := con.getSign()
	paths := "/open/door/openSingleDoor?appid=" + appid + "&timestamp=" + t + "&sign=" + sign + "&doorId=" + strconv.Itoa(id) + "&time=" + timestamp
	return con.Http(paths, "")
}

// 生成签名
func (con ControlServer) getSign() (string, string) {
	t := strconv.Itoa(int(time.Now().Unix())) + "000"
	return ay.MD5(secret + t + appid), t
}

// GetId 获取门禁id
func (con ControlServer) GetId(sn string) int {
	sign, t := con.getSign()
	paths := "/open/device/getBySn?appid=" + appid + "&timestamp=" + t + "&sign=" + sign + "&sn=" + sn
	res := con.Http(paths, "")
	type r struct {
		Ret    int    `json:"ret"`
		ErrMsg string `json:"errMsg"`
		Data   struct {
			CreateTime     int64  `json:"createTime"`
			ModifyTime     int64  `json:"modifyTime"`
			Id             int    `json:"id"`
			Name           string `json:"name"`
			Sn             string `json:"sn"`
			Mac            string `json:"mac"`
			Ip             string `json:"ip"`
			Port           string `json:"port"`
			OrgId          int    `json:"orgId"`
			ModelType      int    `json:"modelType"`
			ModelTypeStr   string `json:"modelTypeStr"`
			Online         bool   `json:"online"`
			LastOnlineTime int64  `json:"lastOnlineTime"`
			AreaId         int    `json:"areaId"`
			AreaName       string `json:"areaName"`
			Remark         string `json:"remark"`
			HasBind        bool   `json:"hasBind"`
			DeviceType     int    `json:"deviceType"`
			DeviceTypeStr  string `json:"deviceTypeStr"`
			NetType        string `json:"netType"`
			Disabled       bool   `json:"disabled"`
			DoorList       []struct {
				CreateTime int64         `json:"createTime"`
				ModifyTime int64         `json:"modifyTime"`
				Id         int           `json:"id"`
				DeviceId   int           `json:"deviceId"`
				OrgId      int           `json:"orgId"`
				Name       string        `json:"name"`
				Nums       int           `json:"nums"`
				DoorType   int           `json:"doorType"`
				DoorStatus int           `json:"doorStatus"`
				CanAttence bool          `json:"canAttence"`
				Remark     string        `json:"remark"`
				Disabled   bool          `json:"disabled"`
				Status     int           `json:"status"`
				Online     bool          `json:"online"`
				PersionVo  []interface{} `json:"persionVo"`
			} `json:"doorList"`
			DeviceFlowList []interface{} `json:"deviceFlowList"`
			Delaye         int           `json:"delaye"`
			DoorListStr    string        `json:"doorListStr"`
			Version        string        `json:"version"`
			RowVersion     string        `json:"rowVersion"`
		} `json:"data"`
		Success bool `json:"success"`
	}

	var rj r
	json.Unmarshal([]byte(res), &rj)

	if rj.Data.DoorList == nil {
		return 0
	}
	return rj.Data.DoorList[0].Id

}

func (con ControlServer) Http(paths string, param string) string {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://www.wacloud.net:9130"+paths, strings.NewReader(string(param)))

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		log.Println(err)
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	return string(body)
}
