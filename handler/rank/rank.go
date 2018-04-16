package rank

import (
	"fmt"
	"strconv"
	"time"

	handlerChallenge "github.com/h3llo-hacker/fc/handler/challenge"
	handlerTemplate "github.com/h3llo-hacker/fc/handler/template"
	handlerUser "github.com/h3llo-hacker/fc/handler/user"
	"github.com/h3llo-hacker/fc/types"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

type userSucceedChallenges map[string]time.Duration

var templateScores map[string]float32

func UpdateUsersRank() error {
	log.Info("Update Users Rank")

	// update templates' success rate and calculate its scroe
	refreshTemplate()

	users, err := getAllValidUser()
	if err != nil {
		return err
	}
	for _, user := range users {
		// DO NOT USE Goroutain
		err = updateUserRank(user.UserID)
		if err != nil {
			log.Errorf("Update user[%v]'s Rank Error:[%v]", user.UserID, err)
		}

	}
	log.Info("Update Users Rank Done.")
	return nil
}

// get all templates
func getAllTemplates() ([]types.Template, error) {
	return handlerTemplate.QueryAllTemplates(0, 0, []string{""})
}

// get all users
func getAllValidUser() ([]types.User, error) {
	items := []string{"UserID"}
	return handlerUser.QueryUserAll(items, 0, 0)
}

func updateUserRank(uid string) error {
	rank := float32(0)
	userSucceedChallenges, err := getUserSucceedChallenges(uid)
	if err != nil {
		log.Errorf("updateUserRank Error:[%v]", err)
		return err
	}
	for templateID, _ := range userSucceedChallenges {
		if templateScores[templateID] != 0.0 {
			rank += templateScores[templateID]
		}
	}
	RANK := fmt.Sprintf("%.2f", rank)
	log.Debugf("User:[%v] Rank:[%v]", uid, RANK)

	// update user's rank
	update := bson.M{"$set": bson.M{"Rank": RANK}}
	user := handlerUser.User{
		UserID: uid,
	}
	return user.UpdateUser(update)
}

func getUserSucceedChallenges(uid string) (userSucceedChallenges, error) {
	u := handlerUser.User{
		UserID: uid,
	}
	challenges, err := u.QueryUserChallenges([]string{"succeeded"}, 0, 0)
	if err != nil {
		return nil, err
	}
	return uniqUserChallenges(challenges), nil
}

func uniqUserChallenges(cs []types.UserChallenge) userSucceedChallenges {
	ucs := make(userSucceedChallenges, 0)
	for _, c := range cs {
		if ucs[c.TemplateID] == 0 {
			ucs[c.TemplateID] = c.FinishTime.Sub(c.CreateTime)
		}
		// } else { // TODO where to judge the time user spend on this challenge
		// 	if ucs[c.TemplateID] > c.FinishTime.Sub(c.CreateTime) {
		// 		ucs[c.TemplateID] = c.FinishTime.Sub(c.CreateTime)
		// 	}
		// }
	}
	return ucs
}

// refresh success rate and its score
func refreshTemplate() error {
	templateScores = make(map[string]float32, 0)

	templates, err := getAllTemplates()
	if err != nil {
		return err
	}
	for _, t := range templates {
		// update template's success rate
		rate := calculateTemplateSuccessRate(t.ID)
		update := bson.M{"$set": bson.M{"SuccRate": rate}}
		err = handlerTemplate.UpdateTemplate(t.ID, update)
		if err != nil {
			log.Errorf("Update Template SuccRate Error: [%v]", err)
		}
	}

	templates, err = getAllTemplates()
	if err != nil {
		return err
	}
	for _, t := range templates {
		// update template's score
		score := calculateTemplateScore(t)
		log.Debugf("Score:[%v]", score)
		update := bson.M{"$set": bson.M{"Score": score}}
		err = handlerTemplate.UpdateTemplate(t.ID, update)
		if err != nil {
			log.Errorf("Update Template Score Error: [%v]", err)
		}
		// inject to this dict
		templateScores[t.ID] = score
	}

	return nil
}

func calculateTemplateScore(template types.Template) float32 {
	level := template.Level
	if level == "" {
		level = "0"
	}
	levelN, err := strconv.Atoi(level)
	if err != nil {
		levelN = 1
		log.Errorf("Rank -> calculateTemplateScore Error: [%v]", err)
	}
	rate := template.SuccRate

	score := (1.0 - rate) * 10.0 * float32(levelN)
	// Fscore := strconv.FormatFloat(float64(score), 'e', 2, 64)
	// score64, _ := strconv.ParseFloat(Fscore, 64)
	// log.Debugf("template:[%v] Score: [%v]", template.ID, float32(score64))
	// return float32(score64)
	return score
}

func calculateTemplateSuccessRate(templateID string) float32 {
	filter := bson.M{"TemplateID": templateID}
	selector := bson.M{"UserID": 1, "State": 1}
	challenges, err := handlerChallenge.QueryChallenges(filter, selector)
	if err != nil {
		log.Errorf("update users rank -> calculate template success rate Error:[%v]", err)
	}

	UniqChallenges := uniqChallenges(challenges)
	success := 0
	total := len(UniqChallenges)
	for _, c := range UniqChallenges {
		if c.State == "succeeded" {
			success += 1
		}
	}
	if total == 0 {
		return 0
	}
	// total + 1 : to prevent there is only one challenge deployed and if he secceed, he will got no rank. That's not fair.
	rate := float32(success) / float32(total+1)
	log.Debugf("template:[%v] success rate: [%v]", templateID, rate)
	return rate
}

// one user can create multi challenges whitch has the same templateID, uniq them.
func uniqChallenges(cs []types.Challenge) map[string]types.Challenge {
	uniqcs := make(map[string]types.Challenge, 0)
	for _, c := range cs {
		key := c.UserID + ":" + c.State
		if uniqcs[key].UserID == "" {
			uniqcs[key] = c
		}
	}
	return uniqcs
}
