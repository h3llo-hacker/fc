package utils

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func GenerateChallengeUrl(userUrl, challengeID string) string {
	log.Debugf("Generate Challenge Url, userUrl: [%v], challengeID: [%v]", userUrl, challengeID)
	urlItems := strings.Split(userUrl, "-")
	p1 := urlItems[0]
	p2 := urlItems[len(urlItems)-1]
	p3 := challengeID[len(challengeID)-8 : len(challengeID)-1]
	url := ""
	if p1 == p2 {
		url = fmt.Sprintf("%s-%s", p1, p3)
	} else {
		url = fmt.Sprintf("%s-%s-%s", p1, p2, p3)
	}
	log.Debugf("ChallengeUrl: [%v]", url)
	return url
}
