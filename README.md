# fc
first controller

Master: [![Master Build Status](https://travis-ci.com/wrfly/fc.svg?token=LqBN16z2mHbvTyyYr9hc&branch=master)](https://travis-ci.com/wrfly/fc)

Develop: [![Develop Build Status](https://travis-ci.com/wrfly/fc.svg?token=LqBN16z2mHbvTyyYr9hc&branch=develop)](https://travis-ci.com/wrfly/fc)

### MongoDB
```bash
docker run -dti -p 27017:27017 -v ./mongodb:/data/db mongo:3.4
use admin
db.createUser(
  {
    user: "myUserAdmin",
    pwd: "abc123",
    roles: [ { role: "userAdminAnyDatabase", db: "admin" } ]
  }
)

#
# after restart the container
#

docker run -dti -p 27017:27017 -v ./mongodb:/data/db mongo:3.4 --auth

use test
db.createUser(
  {
    user: "myTester",
    pwd: "xyz123",
    roles: [ { role: "readWrite", db: "test" }]
  }
)
```
set auth: <https://docs.mongodb.com/manual/tutorial/enable-authentication/>


### etcd
```bash
mkdir -p /data/etcd_data
ETH0_IP=$(ip a s eth0 | grep inet | sed "s/.*inet\ \(.*\)\/.*/\1/g")
docker run -dti --name etcd --network host \
    -v /data/etcd_data:/default.etcd \
    -e ETCD_LISTEN_PEER_URLS=http://${ETH0_IP}:2380 \
    -e ETCD_LISTEN_CLIENT_URLS=http://${ETH0_IP}:2379 \
    -e ETCD_ADVERTISE_CLIENT_URLS=http://${ETH0_IP}:2379 \
    wrfly/etcd:latest


PPP0_IP=$(ip a s ppp0 | grep inet | head -1 | sed "s/.*inet\ \(.*\)\/.*/\1/g")
docker run -dti --name etcd --network host \
    -v /home/mr/test/etcd_data:/default.etcd \
    -e ETCD_LISTEN_PEER_URLS=http://${PPP0_IP}:2380 \
    -e ETCD_LISTEN_CLIENT_URLS=http://${PPP0_IP}:2379 \
    -e ETCD_ADVERTISE_CLIENT_URLS=http://${PPP0_IP}:2379 \
    wrfly/etcd:latest
```

# worker node
```bash
iptables -A INPUT -p icmp --icmp-type 8 -s 172.17.0.0/16 -j DROP
```

# master node
```bash
iptables -A INPUT -p tcp --dport 2375 -s 172.17.0.0/16 -j REJECT
```

# web node
```bash
iptables -A INPUT -p tcp --dport 2379 -s 172.17.0.0/16 -j ACCEPT
iptables -A INPUT -p tcp --dport 2379 -s 0/0 -j REJECT
iptables -A INPUT -p tcp --dport 2380 -s 172.17.0.0/16 -j ACCEPT
iptables -A INPUT -p tcp --dport 2380 -s 0/0 -j REJECT
```

export FC_CONFIG="/home/mr/Documents/work_space/fc/bin/config.json"

