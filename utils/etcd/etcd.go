package etcd

import (
	"time"

	"github.com/h3llo-hacker/fc/config"

	"github.com/coreos/etcd/client"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type API struct {
	K client.KeysAPI
}

var ctx = context.Background()

func KeysAPI(etcdConf config.Etcd_struct) (API, error) {
	var A API
	cfg := client.Config{
		Endpoints: etcdConf.Hosts,
		Transport: client.DefaultTransport,
		Username:  etcdConf.User,
		Password:  etcdConf.Pass,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: 1 * time.Second,
	}

	c, err := client.New(cfg)
	if err != nil {
		log.Error(err)
		return A, err
	}
	A.K = client.NewKeysAPI(c)
	return A, nil
}

func (Kapi API) CreateDir(dir string) error {
	opts := &client.SetOptions{
		Dir: true,
	}
	_, err := Kapi.K.Set(ctx, dir, "", opts)
	if err != nil {
		return err
	}
	return nil
}

func (Kapi API) DeleteDir(dir string) error {
	opts := &client.DeleteOptions{
		Dir:       true,
		Recursive: true,
	}
	_, err := Kapi.K.Delete(ctx, dir, opts)
	if err != nil {
		return err
	}
	return nil
}

func (Kapi API) SetValue(key, value string,
	ttl time.Duration) error {

	// Values <= 0 are ignored.
	opts := &client.SetOptions{
		TTL: ttl,
	}

	_, err := Kapi.K.Set(ctx, key, value, opts)
	if err != nil {
		return err
	}
	return nil
}

func (Kapi API) GetValue(key string) (string, error) {
	resp, err := Kapi.K.Get(ctx, key, nil)
	if err != nil {
		return "", err
	}
	return string(resp.Node.Value), nil
}

func (Kapi API) DeleteKey(key string) error {
	_, err := Kapi.K.Delete(ctx, key, nil)
	if err != nil {
		return err
	}
	return nil
}

func (Kapi API) CreateInOrder(dir string) error {
	_, err := Kapi.K.CreateInOrder(ctx, dir, "", nil)
	if err != nil {
		return err
	}
	resp, err := Kapi.K.Get(ctx, dir, nil)
	if err != nil {
		return err
	}
	log.Infoln(resp.Node.Key)
	return nil
}
