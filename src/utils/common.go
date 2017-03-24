package utils

import (
	"fmt"
	// "strings"

	"github.com/nu7hatch/gouuid"
)

func Guuid() string {
	uid, _ := uuid.NewV4()
	guuid := fmt.Sprintf("%v", uid)
	return guuid
	// r := strings.Replace(guuid, "-", "", -1)
	// return r
}
