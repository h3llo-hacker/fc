package utils

import (
	"fmt"
	"strings"

	"github.com/nu7hatch/gouuid"
)

func Guuid() string {
	uid, _ := uuid.NewV4()
	guuid := fmt.Sprintf("%v", uid)
	return guuid
	// r := strings.Replace(guuid, "-", "", -1)
	// return r
}

func BB64(str string) (nstr string) {
	basic := strings.SplitAfter(str, "=")
	basic_ascii := basic[0][:len(basic[0])-1]
	for i := 0; i < len(basic_ascii); i += 2 {
		if i+2 <= len(basic_ascii) {
			nstr += basic_ascii[i+1 : i+2]
			nstr += basic_ascii[i : i+1]
		}
	}
	if len(basic_ascii)-len(nstr) == 1 {
		nstr += basic_ascii[len(basic_ascii)-1 : len(basic_ascii)]
	}
	for i := 1; i < len(basic); i++ {
		nstr += basic[i]
	}
	nstr += "="
	return nstr
}
