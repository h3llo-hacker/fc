package challenge

import (
	"config"
	"fmt"
	"handler/docker"
	"handler/template"
	U "handler/user"
	"strings"
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
		State:       "creating",
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

	// deploy stack
	endpoint := config.Conf.Endpoints[0]
	stackName := challengeID
	deploylogs, err := docker.DeployStack(endpoint, composeFile,
		stackName)
	if err != nil {
		return challengeID, err
	}
	log.Debugf("Deploy Stack [%v] OK, logs: [%v]",
		stackName, deploylogs)

	// push this challenge to user's challenges
	var user U.User
	user.UserID = userID
	update := bson.M{"$push": bson.M{"Challenges": userChallenge}}
	err = user.UpdateUser(update)
	if err != nil {
		return challengeID, err
	}

	go updateChallenge(userID, challengeID)

	return challengeID, nil
}

func AllChallenges() ([]types.Challenge, error) {
	query := bson.M{}
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
	return challenges, nil
}

func QueryChallenge(challengeID string) (types.Challenge, error) {
	var challenges []types.Challenge
	query := bson.M{"ID": challengeID}

	mongo, dbName, err := db.MongoConn()
	if err != nil {
		return types.Challenge{}, err
	}
	db := mongo.DB(dbName)
	collection := db.C(C)
	err = collection.Find(query).Select(nil).All(&challenges)
	if err != nil {
		return types.Challenge{}, err
	}

	if len(challenges) == 0 {
		return types.Challenge{}, fmt.Errorf("Challenge [%v] Not Found!",
			challengeID)
	}

	return challenges[0], nil
}

func RmChallenge(userID, challengeID string) error {
	challenge, err := QueryChallenge(challengeID)
	if challenge.State == "terminated" {
		return fmt.Errorf("Challenge already removed")
	}

	stackName := challengeID
	endpoint := config.Conf.Endpoints[0]
	err = docker.RemoveStack(endpoint, stackName)
	if err != nil {
		return err
	}

	err = updateUserChallengeState(userID, challengeID, "terminated")
	if err != nil {
		return err
	}

	// update states and finish time in collection `challenges`
	err = UpdateChallengeState(challengeID, "terminated")
	if err != nil {
		return err
	}
	return nil
}

func updateUserChallengeState(uid, cid, state string) error {
	selector := bson.M{"Challenges.ChallengeID": cid}
	update := bson.M{"$set": bson.M{"Challenges.$.State": state}}
	err := db.MongoUpdate("user", selector, update)
	if err != nil {
		return err
	}
	return nil
}

func UpdateChallengeState(challengeID, state string) error {
	selector := bson.M{"ID": challengeID}
	update := bson.M{"$set": bson.M{"Time.FinishTime": time.Now(), "State": state}}
	err := db.MongoUpdate(C, selector, update)
	if err != nil {
		log.Errorf("UpdateChallengeState Error: [%v]", err)
		return err
	}
	return nil
}

func ChallengeExist(challengeID string) bool {
	query := bson.M{"ID": challengeID}
	challenges, err := db.MongoFind(C, query, nil)
	if err != nil {
		return false
	}
	if len(challenges) == 0 {
		return false
	}
	return true
}

func ValidateFlag(flag, challengeID string) bool {
	challenge, err := QueryChallenge(challengeID)
	if err != nil {
		return false
	}
	if challenge.Flag == flag {
		return true
	}
	return false
}

func updateUserChallengeServices(challengeID string,
	services []types.Service_struct) error {
	selector := bson.M{"Challenges.ChallengeID": challengeID}
	update := bson.M{"$set": bson.M{"Challenges.$.Services": services}}
	err := db.MongoUpdate("user", selector, update)
	if err != nil {
		return err
	}
	return nil
}

func updateChallenge(userID, challengeID string) {
	tries := 0
START:
	// wait 20s for startup
	time.Sleep(time.Second * 20)
	// just get the first endpoint, maybe it's a bug
	endpoint := config.Conf.Endpoints[0]
	namespace := challengeID
	tasks, err := docker.PsStack(endpoint, namespace)
	if err != nil {
		log.Error("Ps Task Error: [%v], challengeID: [%v]", err, namespace)
		return
	}

	// get port map
	var (
		ChallengeServices []types.Service_struct
		challengeService  types.Service_struct
	)
	for _, task := range tasks {
		service, err := docker.InspectService(task.ServiceID)
		if err != nil {
			log.Errorf("inspect service [%v] error: [%v]", task.ServiceID, err)
			return
		}
		serviceName := strings.SplitAfter(service.Spec.Name, "_")
		challengeService.ServiceName = serviceName[1]
		for _, ports := range service.Endpoint.Ports {
			challengeService.TargetPort = int(ports.TargetPort)
			challengeService.PublishedPort = int(ports.PublishedPort)
		}
		log.Debugf("challengeService: [%v]", challengeService)
		ChallengeServices = append(ChallengeServices, challengeService)
	}
	log.Debugf("ChallengeServices: [%v]", ChallengeServices)

	// update challenge's port state
	err = updateUserChallengeServices(challengeID, ChallengeServices)
	if err != nil {
		log.Errorf("updateUserChallengeServices Error: [%v]", err)
		return
	}

	running_services := 0
	for _, task := range tasks {
		if task.Status.State == "running" {
			running_services += 1
		}
	}
	if running_services != len(tasks) {
		log.Errorf("Challenge not running, id [%v]", challengeID)
		// restart
		if tries == 4 {
			log.Errorf("Challenge not running for 4 times, id [%v]", challengeID)
			return
		}
		tries += 1
		goto START
	}
	// update DB
	err = UpdateChallengeState(challengeID, "running")
	if err != nil {
		log.Errorf("Challenge.Create.UpdateStates Error! [%v] ", err)
		return
	}
	err = updateUserChallengeState(userID, challengeID, "running")
	if err != nil {
		log.Errorf("Challenge.Create.UpdateUserChallengeStates Error! [%v] ", err)
		return
	}

}
