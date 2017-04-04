package monitor

import (
	"config"
	"handler/challenge"
	"time"
	"types"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

func ScanTimoutChallenges() ([]types.Challenge, error) {
	m := config.Conf.ChallengeDuration * time.Hour
	ago := time.Now().Add(-m)
	early := bson.M{"$lte": ago}
	filter := bson.M{"State": "running", "Time.CreateTime": early}
	selector := bson.M{"UserID": 1, "ID": 1}
	challenges, err := challenge.QueryChallenges(filter, selector)
	if err != nil {
		log.Errorf("query challenges error: [%v]", err)
		return nil, err
	}
	return challenges, nil
}

func RemoveTimeOutChallenges() {
	log.Info("Scanning TimeOut Challenges at: [%v]", time.Now())
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
