package str

import (
	"github.com/json-iterator/go"
)

var j jsoniter.API

func init() {
	j = jsoniter.Config{
		EscapeHTML:             false,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		UseNumber:				true,
	}.Froze()
}

func JsonEncode(i interface{}) (string, error) {
	d, e := j.Marshal(i)
	return string(d), e
}

func JsonDecode(s string, i interface{}) error {
	return j.Unmarshal([]byte(s), i)
}

func JsonEncodeByte(i interface{}) ([]byte, error) {
	d, e := j.Marshal(i)
	return d, e
}

func JsonDecodeByte(s []byte, i interface{}) error {
	return j.Unmarshal(s, i)
}