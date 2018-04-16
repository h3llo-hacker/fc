package register

import (
	"config"
	"testing"
	"types"
)

func Test_Register(t *testing.T) {
	config.LoadConfig()

	service := types.Service{
		ServiceName:   "name",
		PublishedPort: 999,
		TargetPort:    80,
	}
	services := []types.Service{service}
	challengeID := "challengeID"
	challengeUrl := "challengeUrl"
	err := RegisterNewChallenge(challengeID, challengeUrl, services)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("OK")
	}
}

func Test_Unregister(t *testing.T) {
	config.LoadConfig()
	challengeID := "challengeID"
	err := UnregisterChallenge(challengeID)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("OK")
	}
}
