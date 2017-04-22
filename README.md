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

# personal test
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
iptables -I INPUT -p all -s 172.19.0.0/16 -d 172.17.0.0/16 -j REJECT
iptables -I FORWARD -p all -s 172.19.0.0/16 -d 172.17.0.0/16 -j REJECT
```

# master node
```bash
# ping
iptables -A OUTPUT -p icmp -d 172.19.0.0/16 -j DROP
iptables -A OUTPUT -p icmp -d 172.18.0.0/16 -j DROP
iptables -A OUTPUT -p icmp -d 172.17.0.0/16 -j DROP

# docker swarm 2375
iptables -A INPUT -p tcp --dport 2375 -s 172.17.0.0/16 -j DROP
iptables -A INPUT -p udp --dport 2375 -s 172.17.0.0/16 -j DROP

# fc-controller 8080
iptables -A INPUT -p tcp --dport 8080 -s 172.17.128.140 -j ACCEPT
iptables -A INPUT -p tcp --dport 8080 -s 172.17.0.0/16 -j DROP

# sshd 22
iptables -A INPUT -p tcp --dport 22 -s 172.17.128.139 -j ACCEPT
iptables -A INPUT -p tcp --dport 22 -s 172.17.0.0/16 -j DROP
```

# web node
```bash
iptables -A INPUT -p tcp --dport 2379 -s 172.17.0.0/16 -j ACCEPT
iptables -A INPUT -p tcp --dport 2379 -s 0/0 -j REJECT
iptables -A INPUT -p icmp -s 172.19.0.0/16 -j REJECT
```

# log node
```bash
# icmp
iptables -A INPUT -p icmp -s 172.17.0.0/16 -j REJECT

# mongod
iptables -A INPUT -p tcp --dport 27017 -s 172.17.128.137 -j ACCEPT
iptables -A FORWARD -p tcp --dport 27017 -s 172.17.128.137 -j ACCEPT

# sshd
iptables -A INPUT -p tcp --dport 22 -s 172.17.128.139 -j ACCEPT
iptables -A INPUT -p tcp --dport 22 -s 172.17.128.137 -j ACCEPT
iptables -A INPUT -p tcp --dport 22 -s 172.17.0.0/16 -j REJECT

# BAN
iptables -I FORWARD -p tcp --match multiport --dports 27017 -s 172.17.0.0/16 -j REJECT
iptables -I FORWARD -p tcp --match multiport --dports 27017 -s 172.17.128.137 -j ACCEPT
iptables -I FORWARD -p tcp --match multiport --dports 27017 -s 172.17.128.140 -j ACCEPT
```

export FC_CONFIG="/home/mr/Documents/work_space/fc/bin/config.json"

