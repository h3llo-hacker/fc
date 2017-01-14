package handler

import (
	"config"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

func Base64(oriStr string) string {
	s := []byte(oriStr)
	return base64.StdEncoding.EncodeToString(s)
}

func Md5(oriStr string) string {
	e := md5.New()
	e.Write([]byte(oriStr))
	cipherStr := e.Sum(nil)
	encPass := hex.EncodeToString(cipherStr)
	return encPass
}
func Password(oriPass string) string {
	salt := config.Salt
	b64 := Base64(oriPass + salt)
	encPass := Md5(b64)
	return encPass
}
