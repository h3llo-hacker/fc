package register

import (
	"fmt"
	"strconv"

	"github.com/h3llo-hacker/fc/config"
	"github.com/h3llo-hacker/fc/types"
	"github.com/h3llo-hacker/fc/utils/etcd"

	log "github.com/sirupsen/logrus"
)

func RegisterNewChallenge(challengeID, UrlPrefix string,
	services []types.Service) (err error) {

	log.Debugf("Register New Challenge, challengeID: [%v], UrlPrefix: [%v], services: [%v]", challengeID, UrlPrefix, services)

	challengeDir := "/challenges/" + challengeID

	kapi, err := etcd.KeysAPI(config.Conf.Etcd)
	if err != nil {
		return err
	}

	// set /xxx-xxx-xxx-xxx/prefix
	key := fmt.Sprintf("/challenges/%s/prefix", challengeID)
	value := UrlPrefix
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
	// refresh etcd
	refresh := "/challenges/refresh"
	kapi.CreateInOrder(refresh)
	kapi.DeleteDir(refresh)

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
