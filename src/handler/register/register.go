package register

import (
	"config"
	"fmt"
	"strconv"
	"types"
	"utils/etcd"

	log "github.com/Sirupsen/logrus"
)

func RegisterNewChallenge(challengeID, challengeUrl string,
	services []types.Service) (err error) {

	log.Debugf("Register New Challenge, challengeID: [%v], challengeUrl: [%v], services: [%v]", challengeID, challengeUrl, services)

	challengeDir := "/challenges/" + challengeID

	kapi, err := etcd.KeysAPI(config.Conf.Etcd)
	if err != nil {
		return err
	}

	// set /xxx-xxx-xxx-xxx/prefix
	key := fmt.Sprintf("/challenges/%s/prefix", challengeID)
	value := challengeUrl
	err = kapi.SetValue(key, value, 0)
	if err != nil {
		kapi.DeleteDir(challengeDir)
		return err
	}

	// set services (/xxx-xxx-xxx-xxx/services/nginx{pub, tgt})
	for _, service := range services {
		key = fmt.Sprintf("/challenges/%s/services/%s/pub",
			challengeID, service.ServiceName)
		value = strconv.Itoa(service.PublishedPort)
		err = kapi.SetValue(key, value, 0)
		if err != nil {
			kapi.DeleteDir(challengeDir)
			return err
		}

		key = fmt.Sprintf("/challenges/%s/services/%s/tgt",
			challengeID, service.ServiceName)
		value = strconv.Itoa(service.TargetPort)
		err = kapi.SetValue(key, value, 0)
		if err != nil {
			kapi.DeleteDir(challengeDir)
			return err
		}
	}

	return nil
}

func UnregisterChallenge(challengeID string) (err error) {
	log.Debugf("Unregister Challenge, challengeID: [%v]", challengeID)

	challengeDir := "/challenges/" + challengeID

	kapi, err := etcd.KeysAPI(config.Conf.Etcd)
	if err != nil {
		return err
	}

	err = kapi.DeleteDir(challengeDir)
	if err != nil {
		return err
	}
	return nil
}
