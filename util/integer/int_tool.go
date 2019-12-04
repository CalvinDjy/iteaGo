package integer

import (
	"strconv"
)

func Itos(i int) string {
	return strconv.Itoa(i)
}

func I64tos(i int64) string {
	return strconv.FormatInt(i, 10)
}

func IarrToSarr(i []int) []string {
	var s []string
	for _, v := range i {
		s = append(s, strconv.Itoa(v))
	}
	return s
}