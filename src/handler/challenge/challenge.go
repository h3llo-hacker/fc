package challenge

import (
	"config"
	"fmt"
	"handler/docker"
	"handler/register"
	"handler/template"
	U "handler/user"
	"strings"
	"sync"
	"time"
	"types"
	"utils"
	db "utils/db"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types/swarm"
	"gopkg.in/mgo.v2/bson"
)

var (
	C          = "challenges"
	createLock = make(map[string]sync.Locker, 0)
)

func CreateChallenge(userID, templateID, challengeID string) (string, error) {

	if createLock[userID] == nil {
		createLock[userID] = &sync.Mutex{}
	}
	createLock[userID].Lock()
	defer createLock[userID].Unlock()

	var user U.User
	user.UserID = userID

	if challengeID == "" {
		challengeID = utils.Guuid()
	}

	// query user info
	tu, err := user.QueryUser([]string{"UserURL", "Quota"})
	if err != nil {
		return "", err
	}

	// check if user has enough quota
	challenges, err := user.QueryUserChallenges([]string{"running",
		"creating", "created"}, 0, 0)
	if err != nil {
		return "", err
	}
	if len(challenges) >= tu.Quota {
		return "", fmt.Errorf("User has not enough quota to create challenges.")
	}

	// generate something
	flag := utils.RandomFlag()
	now := time.Now()
	UrlPrefix := utils.GenerateChallengeUrl(tu.UserURL, challengeID)

	log.Debugf("Creating challenge [%v], flag is [%v]", challengeID, flag)

	t, err := template.QueryTemplate(templateID)
	if err != nil {
		return challengeID, err
	}

	// check if template is enabled
	if !t.Enable {
		return "", fmt.Errorf("Template [%v] Disabled", templateID)
	}

	// update database first
	challenge := types.Challenge{
		ID:         challengeID,
		Name:       t.Name,
		TemplateID: t.ID,
		Flag:       flag,
		UrlPrefix:  UrlPrefix,
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
		UrlPrefix:   UrlPrefix,
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
			log.Errorf("Update Challenge Services Error: [%v]",
				err)
			RmChallenge(userID, challengeID)
			return
		}
		// register to etcd
		err = register.RegisterNewChallenge(challengeID, UrlPrefix, services)
		if err != nil {
			log.Errorf("RegisterNewChallenge Error: [%v]", err)
			RmChallenge(userID, challengeID)
			return
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

func AllChallenges(limit, skip int) ([]types.Challenge, error) {
	query := bson.M{}
	var challenges []types.Challenge
	mongo, dbName, err := db.MongoConn()
	if err != nil {
		return nil, err
	}
	db := mongo.DB(dbName)
	collection := db.C(C)
	err = collection.Find(query).Limit(limit).Skip(skip).All(&challenges)
	if err != nil {
		return nil, err
	}
	if len(challenges) == 0 {
		return nil, fmt.Errorf("404 challenges not found")
	}
	return challenges, nil
}

func QueryChallenge(filter, selector bson.M) (types.Challenge, error) {
	var challenges []types.Challenge

	mongo, dbName, err := db.MongoConn()
	if err != nil {
		return types.Challenge{}, err
	}
	db := mongo.DB(dbName)
	collection := db.C(C)
	err = collection.Find(filter).Select(selector).All(&challenges)
	if err != nil {
		return types.Challenge{}, err
	}

	if len(challenges) == 0 {
		return types.Challenge{}, fmt.Errorf("Challenge [%v] Not Found!",
			filter["ID"].(string))
	}

	if challenges[0].State == "created" {
		go func() {
			for i := 0; i < 6; i++ {
				err = RefreshChallengeState(filter["ID"].(string))
				if err == nil {
					break
				}
				log.Errorf("RefreshChallengeState Error: [%v], ChallengeID: [%v]", err, filter["ID"].(string))
				time.Sleep(5 * time.Second)
			}
		}()
	}

	return challenges[0], nil
}

func QueryChallenges(filter, selector bson.M) ([]types.Challenge, error) {
	var challenges []types.Challenge

	mongo, dbName, err := db.MongoConn()
	if err != nil {
		return nil, err
	}
	db := mongo.DB(dbName)
	collection := db.C(C)
	err = collection.Find(filter).Select(selector).All(&challenges)
	if err != nil {
		return nil, err
	}

	return challenges, nil
}

func RmChallenge(userID, challengeID string) error {
	filter := bson.M{"ID": challengeID, "UserID": userID}
	selector := bson.M{"State": 1, "UserID": 1}
	challenge, err := QueryChallenge(filter, selector)
	// important!
	if err != nil {
		log.Errorf("RmChallenge Error: [%v]", err)
		return fmt.Errorf("Challenge not belong to user")
	}
	if challenge.State != "running" {
		return fmt.Errorf("Challenge isn't running.")
	}

	log.Debugf("Remove challenge: [%v]", challengeID)

	// important!
	if challenge.UserID != userID {
		log.Debugf("cuid[%v], uid[%v]", challenge.UserID, userID)
		return fmt.Errorf("Challenge not belong to user")
	}

	stackName := challengeID
	endpoint := config.Conf.Endpoint
	err = docker.RemoveStack(endpoint, stackName)
	if err != nil {
		log.Errorf("Remove Challenge Error: [%v], ChallengeID: [%v], UserID: [%v]", err, challengeID, userID)
		return err
	}

	// update states and finish time in collection `challenges` and `users`
	err = UpdateChallengeState(challengeID, "terminated")
	if err != nil {
		log.Errorf("Update Challenge State Error: [%v], ChallengeID: [%v]", err, challengeID)
		return err
	}

	//  Unregister challenge from etcd
	err = register.UnregisterChallenge(challengeID)
	if err != nil {
		log.Errorf("Unregister Challenge Error: [%v], ChallengeID: [%v]", err, challengeID)
		return err
	}

	return nil
}

func UpdateUserChallengeState(uid, cid, state string) error {
	u := bson.M{}
	if state == "running" {
		u = bson.M{"Challenges.$.State": state}
	} else {
		u = bson.M{"Challenges.$.State": state,
			"Challenges.$.FinishTime": time.Now()}
	}
	update := bson.M{"$set": u}
	selector := bson.M{"Challenges.ChallengeID": cid}
	err := db.MongoUpdate("user", selector, update)
	if err != nil {
		return err
	}
	return nil
}

func UpdateChallengeState(challengeID, state string) error {
	u := bson.M{}
	if state == "running" {
		u = bson.M{"State": state}
	} else {
		u = bson.M{"Time.FinishTime": time.Now(), "State": state}
	}
	update := bson.M{"$set": u}
	query := bson.M{"ID": challengeID}
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
	selector := bson.M{"ID": 1}
	challenges, err := db.MongoFind(C, query, selector)
	if err != nil {
		return false
	}
	if len(challenges) == 0 {
		return false
	}
	return true
}

func ValidateFlag(flag, challengeID string) (string, bool) {
	filter := bson.M{"ID": challengeID}
	selector := bson.M{"Flag": 1, "UserID": 1}
	log.Debugf("Validate flag: [%v], ChallengeID: [%v]",
		flag, challengeID)
	challenge, err := QueryChallenge(filter, selector)
	if err != nil {
		return "", false
	}
	if challenge.Flag == flag {
		return challenge.UserID, true
	}
	return "", false
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

func updateChallenge(userID, challengeID string) ([]types.Service,
	error) {
	var (
		ChallengeServices []types.Service
		singleService     types.Service
		tasks             []swarm.Task
		err               error
		service           swarm.Service
	)

	// just get the first endpoint, maybe it's a bug
	endpoint := config.Conf.Endpoint
	namespace := challengeID
	for i := 0; i < 60; i++ {
		tasks, err = docker.PsStack(endpoint, namespace)
		if err == nil {
			break
		}
		log.Errorf("Ps Task Error: [%v], challengeID: [%v]",
			err, namespace)
		time.Sleep(time.Second * 5)
	}

	// get port map
	for _, task := range tasks {
		for i := 0; i < 10; i++ {
			service, err = docker.InspectService(task.ServiceID)
			if err == nil {
				break
			}
			time.Sleep(3 * time.Second)
		}
		if service.Spec.Name == "" {
			return ChallengeServices,
				fmt.Errorf("Inspect service [%v] error: [%v]",
					task.ServiceID, err)
		}

		serviceName := strings.SplitAfterN(service.Spec.Name,
			"_", 2)
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
		ChallengeServices = append(ChallengeServices,
			singleService)
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
		// log.Errorf("Ps Task Error: [%v], challengeID: [%v]", err, namespace)
		return err
	}

	running_services := 0
	for _, task := range tasks {
		if task.Status.State == "running" {
			running_services += 1
		}
	}
	if running_services != len(tasks) {
		log.Debugf("Challenge not running, id [%v]", challengeID)
		return nil
	}

	// update DB
	err = UpdateChallengeState(challengeID, "running")
	if err != nil {
		e := fmt.Errorf("Challenge.RefreshChallengeState.UpdateStates Error: [%v] ", err)
		return e
	}
	return nil
}
