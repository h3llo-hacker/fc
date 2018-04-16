package email

import (
	"testing"

	"github.com/h3llo-hacker/fc/config"
	"github.com/h3llo-hacker/fc/types"
)

func Test_SendVerifyEmail(t *testing.T) {
	config.LoadConfig()
	receiver := types.ValidateEmail{
		EmailAddr: "root@kfd.me",
		UserName:  "root",
		ClickURL:  "http://kfd.me?Verify",
	}

	if err := SendVerifyEmail(receiver); err != nil {
		t.Error("fail")
	}
	t.Log("ok")
}

func Test_SendResetPwdEmail(t *testing.T) {
	config.LoadConfig()
	receiver := types.ValidateEmail{
		EmailAddr: "root@kfd.me",
		UserName:  "root",
		ClickURL:  "http://kfd.me?ResetPwd",
	}

	if err := SendResetPwdEmail(receiver); err != nil {
		t.Error("fail")
	}
	t.Log("ok")
}
