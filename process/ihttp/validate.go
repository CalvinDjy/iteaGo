package ihttp

import (
	"errors"
	"net/http"
	"strings"
)

func Validate(r *http.Request, rules []map[string]string) (map[string]string, error) {
	r.FormValue("") //init by first time

	data := make(map[string]string)
	l := len(rules)
	ch := make(chan map[string]string, l)
	defer close(ch)

	for _, rule := range rules {
		go func(rule map[string]string) {
			if _, ok := rule["key"]; !ok {
				ch <- map[string]string{"key": ""}
				return
			}
			if strings.EqualFold(rule["key"], "") {
				ch <- map[string]string{"key": ""}
				return
			}
			key, value := getValue(rule["key"], r)
			if strings.EqualFold(value, "") {
				if _, ok := rule["default"]; ok {
					value = rule["default"]
				}
			}
			if _, ok := rule["rule"]; !ok {
				ch <- map[string]string{"key": key, "value": value, "err": ""}
				return
			}
			if checkRule(value, rule["rule"]) {
				ch <- map[string]string{"key": key, "value": value, "err": ""}
				return
			}
			if _, ok := rule["msg"]; ok {
				ch <- map[string]string{"key": key, "value": "", "err": rule["msg"]}
				return
			}
			ch <- map[string]string{"key": key, "value": "", "err": "Parameter [" + key + "] validate error"}
			return
		}(rule)
	}

	hasError := false
	var errMsg []string

	for i := 0; i < l; i++ {
		res := <-ch
		if !strings.EqualFold(res["err"], "") {
			hasError = true
			errMsg = append(errMsg, res["err"])
			continue
		}
		if strings.EqualFold(res["key"], "") {
			continue
		}
		data[res["key"]] = res["value"]
	}

	if hasError {
		return nil, errors.New(strings.Join(errMsg, "; "))
	}

	return data, nil
}

func getValue(key string, r *http.Request) (string, string) {
	keyArr := strings.Split(key, "|")
	if len(keyArr) == 1 {
		return keyArr[0], r.FormValue(key)
	}
	if strings.EqualFold(keyArr[1], "header") {
		return keyArr[0], r.Header.Get(keyArr[0])
	}
	return keyArr[0], ""
}

func checkRule(value string, rule string) bool {
	switch rule {
	case "":
		return true
	case "required":
		if strings.EqualFold(value, "") {
			return false
		}
		return true
	default:
		return false
	}
}
