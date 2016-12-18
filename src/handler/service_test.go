package handler

import (
	log "github.com/Sirupsen/logrus"
	"handler"
)

func CreateServiceTest() {
	endpoint := "127.0.0.1:2374"
	serviceName := "nnnginx"
	serviceImage := "nginx"
	err := handler.CreateService(endpoint, serviceName, serviceImage)
	if err == nil {
		log.Info("OK")
	} else {
		log.Error(err)
	}
}
