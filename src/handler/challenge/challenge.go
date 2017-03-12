package challenge

import (
	"config"
	"fmt"
	"handler/docker"
	"handler/register"
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
	var user U.User
	user.UserID = userID
	tu, err := user.QueryUser([]string{"UserURL"})
	if err != nil {
		return "", err
	}

	flag := utils.RandomFlag()
	now := time.Now()
	challengeUrl := utils.GenerateChallengeUrl(tu.UserURL, challengeID)

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
		Url:        challengeUrl,
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
		Url:         challengeUrl,
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
	endpoint := config.Conf.Endpoint
	stackName := challengeID
	deploylogs, err := docker.DeployStack(endpoint, composeFile,
		stackName)
	if err != nil {
		return challengeID, err
	}
	log.Debugf("Deploy Stack [%v] OK, logs: [%v]",
		stackName, deploylogs)

	// push this challenge to user's challenges
	update := bson.M{"$push": bson.M{"Challenges": userChallenge}}
	err = user.UpdateUser(update)
	if err != nil {
		return challengeID, err
	}
	go func() {
		// update services
		services, err := updateChallenge(userID, challengeID)
		if err != nil {
			log.Errorf("Update Challenge Services Error: [%v]", err)
			return
		}
		// register to etcd
		err = register.RegisterNewChallenge(challengeID, challengeUrl, services)
		if err != nil {
			log.Errorf("RegisterNewChallenge Error: [%v]", err)
		}
		// update challenge state (running)
		for i := 0; i < 6; i++ {
			time.Sleep(10 * time.Second)
			err = RefreshChallengeState(challengeID)
			if err == nil {
				break
			}
		}
	}()

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

	if challenges[0].State == "created" {
		go func() {
			err = RefreshChallengeState(challengeID)
			if err != nil {
				log.Errorf("RefreshChallengeState Error: [%v]", err)
			}
		}()
	}

	return challenges[0], nil
}

func RmChallenge(userID, challengeID string) error {
	challenge, err := QueryChallenge(challengeID)
	if challenge.State == "terminated" {
		return fmt.Errorf("Challenge already removed")
	}

	stackName := challengeID
	endpoint := config.Conf.Endpoint
	err = docker.RemoveStack(endpoint, stackName)
	if err != nil {
		return err
	}

	// update states and finish time in collection `challenges` and `users`
	err = UpdateChallengeState(challengeID, "terminated")
	if err != nil {
		return err
	}

	//  Unregister challenge from etcd
	err = register.UnregisterChallenge(challengeID)
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserChallengeState(uid, cid, state string) error {
	selector := bson.M{"Challenges.ChallengeID": cid}
	update := bson.M{"$set": bson.M{"Challenges.$.State": state}}
	err := db.MongoUpdate("user", selector, update)
	if err != nil {
		return err
	}
	return nil
}

func UpdateChallengeState(challengeID, state string) error {
	query := bson.M{"ID": challengeID}
	update := bson.M{"$set": bson.M{"Time.FinishTime": time.Now(), "State": state}}
	err := db.MongoUpdate(C, query, update)
	if err != nil {
		return err
	}
	selector := bson.M{"UserID": 1}
	c, err := db.MongoFind(C, query, selector)
	if err != nil {
		return err
	}
	uid := c[0].(bson.M)["UserID"].(string)
	err = UpdateUserChallengeState(uid, challengeID, state)
	if err != nil {
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
	services []types.Service) error {
	selector := bson.M{"Challenges.ChallengeID": challengeID}
	update := bson.M{"$set": bson.M{"Challenges.$.Services": services}}
	err := db.MongoUpdate("user", selector, update)
	if err != nil {
		return err
	}
	return nil
}

func updateChallenge(userID, challengeID string) ([]types.Service, error) {
	var (
		ChallengeServices []types.Service
		singleService     types.Service
	)

	// just get the first endpoint, maybe it's a bug
	endpoint := config.Conf.Endpoint
	namespace := challengeID
	tasks, err := docker.PsStack(endpoint, namespace)
	if err != nil {
		log.Error("Ps Task Error: [%v], challengeID: [%v]", err, namespace)
		return ChallengeServices, err
	}

	// get port map
	for _, task := range tasks {
		service, err := docker.InspectService(task.ServiceID)
		if err != nil {
			log.Errorf("inspect service [%v] error: [%v]", task.ServiceID, err)
			return ChallengeServices, err
		}
		serviceName := strings.SplitAfter(service.Spec.Name, "_")
		singleService.ServiceName = serviceName[1]
		// some service doesn't have any publish ports
		if len(service.Endpoint.Ports) == 0 {
			singleService.TargetPort = 0
			singleService.PublishedPort = 0
		} else {
			for _, ports := range service.Endpoint.Ports {
				singleService.TargetPort = int(ports.TargetPort)
				singleService.PublishedPort = int(ports.PublishedPort)
			}
		}
		ChallengeServices = append(ChallengeServices, singleService)
	}
	log.Debugf("ChallengeServices: [%v]", ChallengeServices)

	// update challenge's port state
	err = updateUserChallengeServices(challengeID, ChallengeServices)
	if err != nil {
		log.Errorf("updateUserChallengeServices Error: [%v]", err)
		return ChallengeServices, err
	}
	return ChallengeServices, nil
}

func RefreshChallengeState(challengeID string) error {
	// just get the first endpoint, maybe it's a bug
	endpoint := config.Conf.Endpoint
	namespace := challengeID
	tasks, err := docker.PsStack(endpoint, namespace)
	if err != nil {
		log.Error("Ps Task Error: [%v], challengeID: [%v]", err, namespace)
		return err
	}

	running_services := 0
	for _, task := range tasks {
		if task.Status.State == "running" {
			running_services += 1
		}
	}
	if running_services != len(tasks) {
		log.Errorf("Challenge not running, id [%v]", challengeID)
		return fmt.Errorf("Challenge not running, id [%v]", challengeID)
	}

	// update DB
	err = UpdateChallengeState(challengeID, "running")
	if err != nil {
		e := fmt.Errorf("Challenge.RefreshChallengeState.UpdateStates Error: [%v] ", err)
		return e
	}
	return nil
}
