package email

import (
	"testing"
)

func Test_sendMail(t *testing.T) {
	if err := sendMail(); err != nil {
		t.Error("fail")
	}
	t.Log("ok")
}
