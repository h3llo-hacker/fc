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
docker run -dti -p 2379:2379 -p 2380:2380 -v etcd_data:/default.etcd quay.io/coreos/etcd:latest etcd --listen-client-urls http://0.0.0.0:2379 --advertise-client-urls http://0.0.0.0:2380
```

export FC_CONFIG="/home/mr/Documents/work_space/fc/bin/config.json"