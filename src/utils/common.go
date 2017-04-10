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
}

func BB64(str string) string {
	index := strings.Index(str, "=")
	if index == -1 {
		index = len(str)
	}
	equals := str[index:]
	s := ""
	for i := 0; i < index; i += 2 {
		if i+2 <= index {
			s += str[i+1 : i+2]
		}
		s += str[i : i+1]
	}
	return fmt.Sprintf("%s%s", s, equals)
}

func ArrayContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
