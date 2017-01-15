package handler

import (
	"testing"
)

func TestBase64(t *testing.T) {
	str := "hello"
	b64 := "aGVsbG8="
	if Base64(str) == b64 {
		t.Log("Base64 Pass")
	} else {
		t.Error("Base64 Failed")
	}
}

func TestMd5(t *testing.T) {
	oriPass := "hello"
	encPass := "5d41402abc4b2a76b9719d911017c592"
	if Md5(oriPass) == encPass {
		t.Log("Md5 enc Pass.")
	} else {
		t.Error("Md5 enc Faild")
	}
}

func TestPassword(t *testing.T) {
	oriPass := "hello"
	t.Log(Password(oriPass))
}
