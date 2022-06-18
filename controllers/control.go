/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package controllers

import (
	"dragonfly/ay"
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

func (con ControlServer) getSign() (string, string) {
	t := strconv.Itoa(int(time.Now().Unix())) + "000"
	return ay.MD5(secret + t + appid), t
}

// Set 设置设备
func (con ControlServer) Set(sn, status string) string {
	sign, t := con.getSign()
	paths := "/open/door/openSingleDoor?appid=" + appid + "&timestamp=" + t + "&sign=" + sign + "&doorId=1&time=" + status
	return con.Http(paths, "")
}

func (con ControlServer) Http(paths string, param string) string {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://www.wacloud.net:9130"+paths, strings.NewReader(string(param)))
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
