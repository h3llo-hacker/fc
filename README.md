# fc
first controller

Master: [![Master Build Status](https://travis-ci.com/wrfly/fc.svg?token=LqBN16z2mHbvTyyYr9hc&branch=master)](https://travis-ci.com/wrfly/fc)

Develop: [![Develop Build Status](https://travis-ci.com/wrfly/fc.svg?token=LqBN16z2mHbvTyyYr9hc&branch=develop)](https://travis-ci.com/wrfly/fc)

### MongoDB
```bash
docker volume create --name mongodb
docker run -dti -p 27017:27017 -v mongodb:/data/db mongo:3.2 --auth

db.createUser({ user: 'muser', pwd: 'mpass', roles: [ { role: "userAdminAnyDatabase", db: "fc" } ] });
```

### etcd
```bash
docker volume create etcd_data
docker run -ti --network host \
	-v etcd_data:/default.etcd \
	-e ETCD_LISTEN_PEER_URLS=http://10.170.32.166:2380 \
	-e ETCD_LISTEN_CLIENT_URLS=http://10.170.32.166:2379 \
	-e ETCD_ADVERTISE_CLIENT_URLS=http://10.170.32.166:2379 \
	quay.io/coreos/etcd:latest
```

export FC_CONFIG="/home/mr/Documents/work_space/fc/bin/config.json"
