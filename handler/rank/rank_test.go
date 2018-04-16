package rank

import (
	"config"
	"fmt"
	"testing"
)

func Test_getAllTemplates(t *testing.T) {
	config.LoadConfig()

	tps, err := getAllTemplates()
	if err != nil {
		t.Error(err)
	} else {
		for _, tp := range tps {
			fmt.Println(tp.ID)
		}
	}
}

func Test_getAllValidUser(t *testing.T) {
	urs, err := getAllValidUser()
	if err != nil {
		t.Error(err)
	} else {
		for _, ur := range urs {
			fmt.Println(ur.UserID)
		}
	}
}

func Test_getUserSucceedChallenges(t *testing.T) {
	uid := "505a9d03-96d5-4a19-44ae-9a49eb382667"
	usc, err := getUserSucceedChallenges(uid)
	if err != nil {
		t.Error(err)
	}
	for id, c := range usc {
		fmt.Println(id)
		fmt.Println(c)
	}
}

func Test_calculateTemplateSuccessRate(t *testing.T) {
	tid := "fff"
	rate := calculateTemplateSuccessRate(tid)
	fmt.Println(rate)
}

func Test_refreshTemplate(t *testing.T) {
	err := refreshTemplate()
	if err != nil {
		t.Error(err)
	}
}

// updateUserRanks
func Test_updateUserRanks(t *testing.T) {
	err := updateUserRanks("505a9d03-96d5-4a19-44ae-9a49eb382667")
	if err != nil {
		t.Error(err)
	}
}
