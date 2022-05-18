package utils

import "github.com/axgle/mahonia"

// EncodingTo 字符串从 from 编码转到 to 编码
func EncodingTo(s string, from string, to string) string {
	Decoder := mahonia.NewDecoder(from)
	Encoder := mahonia.NewEncoder(to)
	return Encoder.ConvertString(Decoder.ConvertString(s))
}

// GBKToUTF 把 GBK 字符串转成 UTF8 字符串
func GBKToUTF(s string) string {
	return EncodingTo(s, "GBK", "UTF8")
}
