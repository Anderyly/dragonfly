/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package task

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func Start() {
	Coupon()
	Card()
	//Vip()
	Control()
	ControlDel()
	ControlOpenElectric()
	ControlCloseElectric()
	ControlCloseDoor()
	res := Http()

	if res != "1" {
		panic("server error")

	}
}

func Http() string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://blog.aaayun.cc/dragonfly.php", strings.NewReader(string("")))

	//req.Header.Set("Content-Type", "application/json; charset=UTF-8")

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
