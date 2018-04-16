package monitor

import (
	"config"
	"testing"
)

func Test_ScanTimoutChallenges(t *testing.T) {
	config.LoadConfig()
	cs, err := ScanTimoutChallenges()
	if err != nil {
		t.Error(err)
	}
	for _, c := range cs {
		t.Log(c.Time.CreateTime)
	}
}
