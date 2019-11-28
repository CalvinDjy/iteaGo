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
