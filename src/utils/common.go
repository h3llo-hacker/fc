package utils

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
)

func Guuid() string {
	uid, _ := uuid.NewV4()
	return fmt.Sprintf("%v", uid)
}
