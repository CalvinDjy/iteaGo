package str

import (
	"github.com/CalvinDjy/iteaGo/ilog"
	"strconv"
)


func Toint(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		ilog.Error("str to int error : ", err)
	}
	return v
}

func Toint64(s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		ilog.Error("str to int64 error : ", err)
	}
	return v
}

func Tobool(s string) bool {
	v, err := strconv.ParseBool(s)
	if err != nil {
		ilog.Error("str to bool error : ", err)
	}
	return v
}