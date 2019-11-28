package algorithm

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/CalvinDjy/iteaGo/ilog"
	"io"
	"math/rand"
	"time"
)

func Rand(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}

func Md5(s string) string {
	w := md5.New()
	_, err := io.WriteString(w, s)
	if err != nil {
		ilog.Error("md5 error : ", err)
		return ""
	}
	return fmt.Sprintf("%x", w.Sum(nil))
}

func Sha1(s string) string {
	h := sha1.New()
	_, err := io.WriteString(h, s)
	if err != nil {
		ilog.Error("sha1 error : ", err)
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func HashSha1(s string, key string, raw bool) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(s))
	if raw {
		return string(mac.Sum(nil))
	} else {
		return fmt.Sprintf("%x", mac.Sum(nil))
	}

}

func Base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func Base64Decode(s string) string {
	decodeBytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		ilog.Error("base64decode error : ", err)
		return ""
	}
	return string(decodeBytes)
}

func Base64UrlEncode(s string) string {
	return base64.URLEncoding.EncodeToString([]byte(s))
}

func Base64UrlDecode(s string) string {
	decodeBytes, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		ilog.Error("base64urldecode error : ", err)
		return ""
	}
	return string(decodeBytes)
}