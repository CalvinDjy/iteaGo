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
	}.Froze()
}

func JsonEncode(i interface{}) (string, error) {
	d, e := j.Marshal(i)
	return string(d), e
}

func JsonDecode(s string, i interface{}) error {
	return j.Unmarshal([]byte(s), i)
}