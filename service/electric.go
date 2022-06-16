/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"dragonfly/ay"
	"encoding/base64"
	"encoding/hex"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ElectricServer struct {
}

var (
	appKey    = "24862778"
	appSecret = "36d3e0833493fe9839be36a6411133c8"
)

func (con ElectricServer) GetSign(method string, headers map[string]string, url string) string {
	h := "X-Ca-Key:" + headers["X-Ca-Key"] + "\nX-Ca-Nonce:" + headers["X-Ca-Nonce"] + "\nX-Ca-Timestamp:" + headers["X-Ca-Timestamp"] + "\n"
	sign := method + "\n" + headers["Accept"] + "\n\n" + headers["Content-Type"] + "\n\n" + h + url
	hmac := con.genHMAC256([]byte(sign), []byte(appSecret))
	return base64.StdEncoding.EncodeToString(hmac)
}

func (con ElectricServer) SHA256(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	res := hex.EncodeToString(hash.Sum(nil))
	return res
}

func (con ElectricServer) genHMAC256(ciphertext, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(ciphertext))
	hmac := mac.Sum(nil)
	return hmac
}

func (con ElectricServer) Header(paths string, param string) {
	uuid := ay.MD5(uuid.NewV4().String())
	t := strconv.FormatInt(time.Now().Unix(), 10) + "000"
	client := &http.Client{}

	req, err := http.NewRequest("POST", "https://api.link.aliyun.com"+paths, strings.NewReader(string(param)))
	log.Println("https://api.link.aliyun.com" + paths)
	if err != nil {
		log.Println(err)
	}
	headers := make(map[string]string)
	headers["Accept"] = "application/text; charset=UTF-8"
	headers["Content-Type"] = "application/text; charset=UTF-8"
	headers["X-Ca-Key"] = appKey
	headers["X-Ca-Nonce"] = uuid
	headers["X-Ca-Timestamp"] = t
	//headers["X-Ca-Stage"] = "RELEASE"
	headers["X-Ca-Signature-Headers"] = "X-Ca-Key,X-Ca-Nonce,X-Ca-Timestamp"
	headers["X-Ca-Signature"] = con.GetSign("POST", headers, paths)

	for key, header := range headers {
		req.Header.Set(key, header)
	}

	//log.Println(headers)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(body))
	//return string(body), nil
}
