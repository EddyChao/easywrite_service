package util

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"github.com/dlclark/regexp2"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func IsPasswordLegal(password string) bool {
	//密码至少包含 数字和英文，长度8-20
	rgx := regexp2.MustCompile("^(?![\\x21-\\x2F\\x3A-\\x40\\x5B-\\x60\\x7B-\\x7E]+$)(?![\\d]+$)(?![a-zA-Z]+$)[\\x21-\\x7e]{8,20}$", 0)
	m, _ := rgx.MatchString(password)
	return m
}

func Sha256Sum(text string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(text)))
}

func Sha1Sum(text string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(text)))
}

func IsEmail(email string) bool {
	//pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func IsPhone(phone string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(phone)
}

func GetRandomString(n int) string {
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	var result []byte
	for i := 0; i < n; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

func GetRandomCode(n int) string {
	str := "0123456789"
	bytes := []byte(str)
	var result []byte
	for i := 0; i < n; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

func GetKey(args ...string) string {
	if len(args) < 2 {
		log.Fatalln("the len(args) must >= 2")
	}
	return strings.Join(args, "-")
}

func GetTimeTick64() int64 {
	return time.Now().UnixNano() / 1e6
}

func GetTimeTick32() int32 {
	return int32(time.Now().Unix())
}

func GetFormatTime(time time.Time) string {
	return time.Format("20060102")
}

// 基础做法 日期20191025时间戳1571987125435+3位随机数
func GenerateOrderNumber() string {
	date := GetFormatTime(time.Now())
	r := rand.Intn(1000)
	code := fmt.Sprintf("%s%d%03d", date, GetTimeTick64(), r)
	return code
}
