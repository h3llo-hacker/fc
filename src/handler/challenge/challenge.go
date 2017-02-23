package challenge

import (
	"config"
	"fmt"
	"handler/docker"
	"handler/template"
	U "handler/user"
	"time"
	"types"
	"utils"
	db "utils/db"

	log "github.com/Sirupsen/logrus"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/mgo.v2/bson"
)

const (
	C = "challenges"
)

func CreateChallenge(userID, templateID, challengeID string) (string, error) {
	if challengeID == "" {
		uid, _ := uuid.NewV4()
		challengeID = fmt.Sprintf("%v", uid)
	}
	flag := utils.RandomFlag()
	now := time.Now()

	log.Debugf("Creating challenge [%v], flag is [%v]", challengeID, flag)

	ts, err := template.QueryTemplate(templateID)
	if err != nil {
		return challengeID, err
	}
	t := ts[0]

	// update database first
	challenge := types.Challenge{
		ID:         challengeID,
		Name:       t.Name,
		TemplateID: t.ID,
		Flag:       flag,
		StackID:    challengeID,
		UserID:     userID,
		Time: types.Time_struct{
			CreateTime: now,
		},
		State: "creating",
	}
	userChallenge := types.UserChallenge{
		ChallengeID: challengeID,
		TemplateID:  t.ID,
		Flag:        flag,
		CreateTime:  now,
	}
	// add challenge to challenge collections (to overview)
	err = db.MongoInsert(C, challenge)
	if err != nil {
		return challengeID, err
	}

	// generate a compose file with the given flag
	composeFile, err := template.GenerateComposeFile(templateID, flag)
	if err != nil {
		return challengeID, err
	}
	log.Debugln(composeFile)
	endpoint := config.Conf.Endpoints[0]
	stackName := challengeID
	deploylogs, err := docker.DeployStack(endpoint, composeFile,
		stackName)
	if err != nil {
		return challengeID, err
	}
	log.Debugf("Deploy Stack [%v] OK, logs: [%v]",
		stackName, deploylogs)

	// update user
	var user U.User
	user.UserID = userID
	update := bson.M{"$push": bson.M{"Challenges.InProcess": userChallenge}}
	err = user.UpdateUser(update)
	if err != nil {
		return challengeID, err
	}
	return challengeID, nil
}

func QueryChallenges(challengeID string) ([]types.Challenge, error) {
	query := bson.M{"ID": challengeID}
	if challengeID == "all" {
		query = bson.M{}
	}

	var challenges []types.Challenge
	mongo, dbName, err := db.MongoConn()
	if err != nil {
		return nil, err
	}
	db := mongo.DB(dbName)
	collection := db.C(C)
	err = collection.Find(query).Select(nil).All(&challenges)
	if err != nil {
		return nil, err
	}

	if len(challenges) == 0 {
		return nil, fmt.Errorf("Challenge [%v] Not Found!",
			challengeID)
	}

	return challenges, nil
}

func RmChallenge(userID, challengeID string) error {
	stackName := challengeID
	endpoint := config.Conf.Endpoints[0]
	err := docker.RemoveStack(endpoint, stackName)
	if err != nil {
		return err
	}

	user := U.User{
		UserID: userID,
	}

	FinishChallenge, err := getInprocessChallenges(user, challengeID)
	if err != nil {
		return err
	}

	err = pushToFinishChallenges(user, FinishChallenge)
	if err != nil {
		return err
	}

	// remove this challenge from user's inprocess challenge
	update := bson.M{"$pull": bson.M{"Challenges.InProcess": bson.M{"ChallengeID": challengeID}}}
	err = user.UpdateUser(update)
	if err != nil {
		return err
	}

	// update finish time in collection [challenges]
	selector := bson.M{"ID": challengeID}
	update = bson.M{"$set": bson.M{"Time.FinishTime": time.Now()}}
	err = db.MongoUpdate(C, selector, update)
	if err != nil {
		return err
	}
	return nil
}

func pushToFinishChallenges(user U.User, FinishChallenge types.UserChallenge) (err error) {
	update := bson.M{"$push": bson.M{"Challenges.Finished": FinishChallenge}}
	err = user.UpdateUser(update)
	if err != nil {
		return err
	}
	return nil
}

func getInprocessChallenges(user U.User, challengeID string) (types.UserChallenge, error) {
	items := []string{"Challenges.InProcess"}
	data, err := user.QueryUser(items)
	if err != nil {
		return types.UserChallenge{}, err
	}
	InProcessChallenges := data[0].Challenges.InProcess
	for _, c := range InProcessChallenges {
		if c.ChallengeID == challengeID {
			c.FinishTime = time.Now()
			return c, nil
		}
	}
	return types.UserChallenge{}, nil
}

func UpdateChallengeState(challengeID, state string) {
	selector := bson.M{"ID": challengeID}
	update := bson.M{"$set": bson.M{"State": state}}
	err := db.MongoUpdate(C, selector, update)
	if err != nil {
		log.Errorf("UpdateChallengeState Error: [%v]", err)
	}
}
