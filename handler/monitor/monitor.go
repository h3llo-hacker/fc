package monitor

import (
	"time"

	"github.com/h3llo-hacker/fc/config"
	"github.com/h3llo-hacker/fc/handler/challenge"
	"github.com/h3llo-hacker/fc/types"

	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

func ScanTimoutChallenges() ([]types.Challenge, error) {
	m := config.Conf.ChallengeDuration * time.Hour
	ago := time.Now().Add(-m)
	early := bson.M{"$lte": ago}
	state := []bson.M{bson.M{"State": "running"}, bson.M{"State": "created"}}
	filter := bson.M{"$or": state, "Time.CreateTime": early}
	selector := bson.M{"UserID": 1, "ID": 1}
	challenges, err := challenge.QueryChallenges(filter, selector)
	if err != nil {
		log.Errorf("query challenges error: [%v]", err)
		return nil, err
	}
	return challenges, nil
}

func RemoveTimeOutChallenges() {
	log.Infof("Scanning TimeOut Challenges at: [%v]", time.Now())
	TOChallenges, err := ScanTimoutChallenges()
	if err != nil {
		log.Errorf("ScanTimeOutChallenge Error: [%v]", err)
		return
	}
	for _, c := range TOChallenges {
		log.Infof("Remove Time Out Challenge: [%v]", c.ID)
		err := challenge.RmChallenge(c.UserID, c.ID)
		if err != nil {
			log.Errorf("RemoveTimeOutChallenges Error: [%v]", err)
		}
	}
}
