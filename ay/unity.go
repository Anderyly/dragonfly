/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package ay

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/w3liu/go-common/constant/timeformat"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

func LastTime(t int) (msg string) {
	s := (int(time.Now().Unix()) - t) / 60

	switch {
	case s < 60:
		msg = strconv.Itoa(s) + "分钟前"

	case s >= 60 && s < (60*24):
		msg = strconv.Itoa(s/60) + "小时前"
	case s >= (60*24) && s < (60*24*3):
		msg = strconv.Itoa(s/24/60) + "天前"

	default:
		msg = time.Unix(int64(t), 0).Format("2006-01-02 15:04:05")
	}
	return
}

func LastTime1(t int) (msg string) {
	s := (int(time.Now().Unix()) - t) / 60

	switch {
	case s < 60:
		msg = strconv.Itoa(s) + "分钟前"

	case s >= 60 && s < (60*24):
		msg = strconv.Itoa(s/60) + "小时前"
	case s >= (60*24) && s < (60*24*3):
		msg = strconv.Itoa(s/24/60) + "天前"

	default:
		msg = time.Unix(int64(t), 0).Format("2006/01/02")
	}
	return
}

func Int32ToString(n int32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}

var num int64

func MakeOrder(t time.Time) string {
	s := t.Format(timeformat.Continuity)
	m := t.UnixNano()/1e6 - t.UnixNano()/1e9*1e3
	ms := sup(m, 3)
	p := os.Getpid() % 1000
	ps := sup(int64(p), 3)
	i := atomic.AddInt64(&num, 1)
	r := i % 10000
	rs := sup(r, 4)
	n := fmt.Sprintf("%s%s%s%s", s, ms, ps, rs)
	return n
}

func GetRandomString(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)

}

//对长度不足n的数字前面补0
func sup(i int64, n int) string {
	m := fmt.Sprintf("%d", i)
	for len(m) < n {
		m = fmt.Sprintf("0%s", m)
	}
	return m
}

func MD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func AuthCode(str, operation, key string, expiry int64) string {
	// 动态密匙长度，相同的明文会生成不同密文就是依靠动态密匙
	// 加入随机密钥，可以令密文无任何规律，即便是原文和密钥完全相同，加密结果也会每次不同，增大破解难度。
	// 取值越大，密文变动规律越大，密文变化 = 16 的 cKeyLength 次方
	// 当此值为 0 时，则不产生随机密钥
	cKeyLength := 1
	if len(str) < cKeyLength {
		return ""
	}

	if len(str) < 11 && operation == "DECODE" {
		return ""
	}
	// 密匙
	if key == "" {
		key = "#@!^5ebcQJx2Lz6GmcsqNiNHW^!@#"
	}
	key = MD5(key)

	// 密匙a会参与加解密
	keyA := MD5(key[:16])
	// 密匙b会用来做数据完整性验证
	keyB := MD5(key[16:])
	// 密匙c用于变化生成的密文
	keyC := ""
	if operation == "DECODE" {
		keyC = str[:cKeyLength]
	} else {
		sTime := MD5(time.Now().String())
		sLen := 32 - cKeyLength
		keyC = sTime[sLen:]
	}
	// 参与运算的密匙
	cryptKey := fmt.Sprintf("%s%s", keyA, MD5(keyA+keyC))
	keyLength := len(cryptKey)
	// 明文，前10位用来保存时间戳，解密时验证数据有效性，10到26位用来保存$keyB(密匙b)，解密时会通过这个密匙验证数据完整性
	// 如果是解码的话，会从第$ckey_length位开始，因为密文前$ckey_length位保存 动态密匙，以保证解密正确
	if operation == "DECODE" {
		str = strings.Replace(str, "-", "+", -1)
		str = strings.Replace(str, "_", "/", -1)
		strByte, err := base64.StdEncoding.DecodeString(str[cKeyLength:])
		if err != nil {
			return ""
		}
		str = string(strByte)
	} else {
		if expiry != 0 {
			expiry = expiry + time.Now().Unix()
		}
		tmpMd5 := MD5(str + keyB)
		str = fmt.Sprintf("%010d%s%s", expiry, tmpMd5[:16], str)
	}
	string_length := len(str)
	resdata := make([]byte, 0, string_length)
	var rndkey, box [256]int
	// 产生密匙簿
	j := 0
	a := 0
	i := 0
	tmp := 0
	for i = 0; i < 256; i++ {
		rndkey[i] = int(cryptKey[i%keyLength])
		box[i] = i
	}
	// 用固定的算法，打乱密匙簿，增加随机性，好像很复杂，实际上并不会增加密文的强度
	for i = 0; i < 256; i++ {
		j = (j + box[i] + rndkey[i]) % 256
		tmp = box[i]
		box[i] = box[j]
		box[j] = tmp
	}
	// 核心加解密部分
	a = 0
	j = 0
	tmp = 0
	for i = 0; i < string_length; i++ {
		a = (a + 1) % 256
		j = (j + box[a]) % 256
		tmp = box[a]
		box[a] = box[j]
		box[j] = tmp
		// 从密匙簿得出密匙进行异或，再转成字符
		resdata = append(resdata, byte(int(str[i])^box[(box[a]+box[j])%256]))
	}
	result := string(resdata)
	if operation == "DECODE" {
		// substr($result, 0, 10) == 0 验证数据有效性
		// substr($result, 0, 10) - time() > 0 验证数据有效性
		// substr($result, 10, 16) == substr(md5(substr($result, 26).$keyB), 0, 16) 验证数据完整性
		// 验证数据有效性，请看未加密明文的格式
		frontTen, _err := strconv.ParseInt(result[:10], 10, 0)
		if _err != nil {
			return ""
		}
		if (frontTen == 0 || frontTen-time.Now().Unix() > 0) && result[10:26] == MD5(result[26:] + keyB)[:16] {
			return result[26:]
		} else {
			return ""
		}
	} else {
		// 把动态密匙保存在密文里，这也是为什么同样的明文，生产不同密文后能解密的原因
		// 因为加密后的密文可能是一些特殊字符，复制过程可能会丢失，所以用base64编码
		result = keyC + base64.StdEncoding.EncodeToString([]byte(result))
		result = strings.Replace(result, "+", "-", -1)
		result = strings.Replace(result, "/", "_", -1)
		return result
	}
}

func MakeCoupon(coupon string) float64 {
	couponArr := strings.Split(coupon, "-")

	log.Println(couponArr)

	gj, err := strconv.ParseFloat(couponArr[1], 64)
	if err != nil {
		log.Println("优惠价错误")
	}
	dj, err1 := strconv.ParseFloat(couponArr[0], 64)
	if err1 != nil {
		log.Println("优惠价错误")
	}
	cha := gj - dj
	var price float64
	for {

		p := 0.01 + rand.Float64()*(cha-0.01)
		price, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", p), 64)
		if price <= cha {
			break
		}
	}
	return price + dj
}

func IsFileExist(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true

}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateMutiDir(filePath string) error {
	if !IsFileExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			log.Println("创建文件夹失败,error info:", err)
			return err
		}
		return err
	}
	return nil
}

func WriteFile(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

func Summary(content string, count int) string {

	content = strings.Replace(content, `@<script(.*?)</script>@is`, "", -1)
	content = strings.Replace(content, `@<iframe(.*?)</iframe>@is`, "", -1)
	content = strings.Replace(content, `@<style(.*?)</style>@is`, "", -1)
	content = strings.Replace(content, `\`, "", -1)

	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	content = re.ReplaceAllString(content, "")

	content = strings.Replace(content, " ", "", -1)
	content = strings.Replace(content, "　　", "", -1)
	content = strings.Replace(content, "\t", "", -1)
	content = strings.Replace(content, "\r", "", -1)
	content = strings.Replace(content, "\n", "", -1)
	cont := []rune(content)

	return string(cont[:count])
}

func GetDirAll(localDir string) int {
	var num int
	filepath.Walk(localDir, func(filename string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}
		num++
		return nil
	})
	return num
}
