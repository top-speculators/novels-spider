package helpers

import "golang.org/x/text/encoding/simplifiedchinese"

// GBK 转 UTF8
func GBKToUTF8(html string) (str string, err error) {
	var str1 []byte
	str1, err = simplifiedchinese.GB18030.NewDecoder().Bytes([]byte(html))
	str = string(str1)
	return
}

// UTF8 转 GBK
func UTF8ToGBK(html string) (str string, err error) {
	var str1 []byte
	str1, err = simplifiedchinese.GB18030.NewEncoder().Bytes([]byte(html))
	str = string(str1)
	return
}
