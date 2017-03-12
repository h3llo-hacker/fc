package utils

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func GenerateChallengeUrl(userUrl, challengeID string) string {
	log.Debugf("Generate Challenge Url, userUrl: [%v], challengeID: [%v]", userUrl, challengeID)
	urlItems := strings.Split(userUrl, "-")
	challengeItems := strings.Split(challengeID, "-")
	url := fmt.Sprintf("%s-%s-%s", urlItems[0], urlItems[len(urlItems)-1], challengeItems[len(challengeItems)-1])
	log.Debugf("ChallengeUrl: [%v]", url)
	return url
}
