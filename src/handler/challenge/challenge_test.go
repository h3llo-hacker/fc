package challenge

import (
	"testing"
)

func Test_CreateChallenge(t *testing.T) {
	userID := "537326cd-a113-4dbf-49f9-93234ec8799a"
	templateID := "5ba174a1-cb81-4227-5f65-2a6c7985f6ea"
	_, err := CreateChallenge(userID, templateID, "")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Log("challenge created")
}

func Test_RmChallenge(t *testing.T) {
	userID := "537326cd-a113-4dbf-49f9-93234ec8799a"
	templateID := "5ba174a1-cb81-4227-5f65-2a6c7985f6ea"
	challengeID, _ := CreateChallenge(userID, templateID, "")
	err := RmChallenge(userID, challengeID)
	if err != nil {
		t.Error(err)
	}
}
