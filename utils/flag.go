package utils

import (
	"fmt"

	"github.com/Pallinder/go-randomdata"
)

func RandomFlag() string {
	s1 := randomdata.Title(randomdata.RandomGender)
	s2 := randomdata.LastName()
	s3 := randomdata.Month()
	s4 := randomdata.Day()
	flag := fmt.Sprintf("flag{%v-%v-%v-%v}", s1, s2, s3, s4)
	return flag
}
