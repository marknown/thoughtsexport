package utils

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// ToJson return a interface's json byte
func ToJson(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// JsonpCallback return a jsonp format callback
func JsonpCallback() string {
	var cstZone = time.FixedZone("CST", 8*3600) // UTC/GMT +08:00
	t := time.Now()
	return "jQuery" + t.In(cstZone).Format("20060102150405")
}

// JSONToStruct 把 Json 字符串解析成指定结构体
func JSONToStruct(s string, v interface{}) error {
	return json.Unmarshal([]byte(s), v)
}

// JSONPToStruct 把 Jsonp 字符串解析成指定结构体
func JSONPToStruct(s string, v interface{}) error {
	pos1 := strings.Index(s, "(")
	pos2 := strings.LastIndex(s, ")")

	if -1 < pos1 && -1 < pos2 {
		s = s[pos1+1 : pos2]
		return JSONToStruct(s, v)
	}

	return errors.New("input is not valid jsonp string")
}
