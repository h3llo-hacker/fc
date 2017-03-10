## `GET    /ping`

```bash
{
    "ping": "pong"
}
```

## User
### `GET    /users`
```bash
[
    {
        "EmailAddress": "test@test.com", 
        "UserName": "名字", 
        "UserNum": 0, 
        "UserURL": "ming-zi"
    }, 
    {
        "EmailAddress": "test2@test.com", 
        "UserName": "名字2", 
        "UserNum": 1, 
        "UserURL": "ming-zi-2"
    }, 
    {
        "EmailAddress": "test3@test.com", 
        "UserName": "名字3", 
        "UserNum": 2, 
        "UserURL": "ming-zi-3"
    }
]
```
### `POST   /user/login`
```bash
➜  ~ http --form POST 127.0.0.1:8080/user/login password=b email=test3@test.com
HTTP/1.1 200 OK
Content-Length: 17
Content-Type: application/json; charset=utf-8
Date: Wed, 08 Mar 2017 08:02:29 GMT

{
    "login": "true"
}

➜  ~ http --form POST 127.0.0.1:8080/user/login password=c email=test3@test.com
HTTP/1.1 401 Unauthorized
Content-Length: 18
Content-Type: application/json; charset=utf-8
Date: Wed, 08 Mar 2017 08:03:23 GMT

{
    "login": "false"
}
```

### `POST   /user/create`

```bash
➜  ~ http --form POST 127.1:8080/user/create username=名字3 password=b ip=223.5.5.5 os=linux ua=chrome email=test4@test.com
HTTP/1.1 200 OK
Content-Length: 18
Content-Type: application/json; charset=utf-8
Date: Wed, 08 Mar 2017 08:05:04 GMT

{
    "Add user": "OK"
}

➜  ~ http --form POST 127.1:8080/user/create username=名字3 password=b ip=223.5.5.5 os=linux ua=chrome email=test4@test.com
HTTP/1.1 500 Internal Server Error
Content-Length: 44
Content-Type: application/json; charset=utf-8
Date: Wed, 08 Mar 2017 08:05:23 GMT

{
    "error": "Email Address Has Already Used."
}

```

### `DELETE /user/delete`
```bash
➜  ~ curl -X DELETE -F "email=test1@test.com" "http://127.0.0.1:8080/user/delete"
{"err":"Remove User Error: [User Email [test1@test.com] Not Found]"}

➜  ~ curl -X DELETE -F "email=test2@test.com" "http://127.0.0.1:8080/user/delete"
{"err":"Remove User Error: [User Email [test2@test.com] Not Found]"}

➜  ~ curl -X DELETE -F "email=test3@test.com" "http://127.0.0.1:8080/user/delete"
{"Rm User OK":"OK"}
➜  ~ 
```

### `POST   /user/update/:userURL`
(TOTO)
```bash
➜  ~ http --form POST 127.1:8080/user/update/ming-zi password=newPassword
HTTP/1.1 200 OK
Content-Length: 18
Content-Type: application/json; charset=utf-8
Date: Wed, 08 Mar 2017 08:05:04 GMT

{
    "update": "OK"
}

```


### `POST   /user/follow/:userURL`
(TODO)
```bash
curl -X POST 127.0.0.1:8080/user/follow/ming-zi -F "user=follow_user_url"
{"follow":"ok"}
# user[ming-zi] follow user[follow_user_url]
```

### `GET    /user/:userURL` & `GET    /user/:userURL/info`
```bash
➜  ~ http 127.0.0.1:8080/user/ming-zi/info
HTTP/1.1 200 OK
Content-Length: 97
Content-Type: application/json; charset=utf-8
Date: Wed, 08 Mar 2017 08:30:02 GMT

{
    "EmailAddress": "test@test.com", 
    "Intro": "", 
    "UserName": "名字", 
    "UserURL": "ming-zi", 
    "WebSite": ""
}

➜  ~ http 127.0.0.1:8080/user/ming-zi     
HTTP/1.1 200 OK
Content-Length: 97
Content-Type: application/json; charset=utf-8
Date: Wed, 08 Mar 2017 08:36:50 GMT

{
    "EmailAddress": "test@test.com", 
    "Intro": "", 
    "UserName": "名字", 
    "UserURL": "ming-zi", 
    "WebSite": ""
}
```

### `GET    /user/:userURL/challenges`
```bash
# all challenges
➜  ~ http 127.0.0.1:8080/user/ming-zi/challenges
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Wed, 08 Mar 2017 08:41:25 GMT
Transfer-Encoding: chunked

[
    {
        "ChallengeID": "0e1e4767-6df9-4869-7936-5eed4db95dde", 
        "CreateTime": "2017-02-24T20:19:52.966+08:00", 
        "FinishTime": "0001-01-01T00:00:00Z", 
        "Flag": "flag{Jayden-Thompson-December-Sunday}", 
        "Services": null, 
        "State": "failed", 
        "Services": [
            {
                "PublishedPort": 0, 
                "ServiceName": "busybox", 
                "TargetPort": 0
            }, 
            {
                "PublishedPort": 30000, 
                "ServiceName": "nginx", 
                "TargetPort": 80
            }
        ], 
        "TemplateID": "5ba174a1-cb81-4227-5f65-2a6c7985f6ea"
    },
    {
        "ChallengeID": "aed54d5e-8e77-4810-4e8e-e078a37e894c", 
        "CreateTime": "2017-03-01T21:18:39.638+08:00", 
        "FinishTime": "0001-01-01T00:00:00Z", 
        "Flag": "flag{Daniel-Smith-July-Monday}", 
        "Services": null, 
        "State": "terminated", 
        "Services": [
            {
                "PublishedPort": 0, 
                "ServiceName": "busybox", 
                "TargetPort": 0
            }, 
            {
                "PublishedPort": 30001, 
                "ServiceName": "nginx", 
                "TargetPort": 80
            }
        ], 
        "TemplateID": "5ba174a1-cb81-4227-5f65-2a6c7985f6ea"
    }, 
    {
        "ChallengeID": "671e4166-db2e-416f-7e33-898ea44e07c2", 
        "CreateTime": "2017-03-01T21:51:53.361+08:00", 
        "FinishTime": "0001-01-01T00:00:00Z", 
        "Flag": "flag{Sofia-Moore-October-Monday}", 
        "Services": [
            {
                "PublishedPort": 0, 
                "ServiceName": "busybox", 
                "TargetPort": 0
            }, 
            {
                "PublishedPort": 30000, 
                "ServiceName": "nginx", 
                "TargetPort": 80
            }
        ], 
        "State": "running", 
        "TemplateID": "5ba174a1-cb81-4227-5f65-2a6c7985f6ea"
    }
]

# type=0 : failed challenges
➜  ~ http "127.0.0.1:8080/user/ming-zi/challenges?type=0"
HTTP/1.1 200 OK
Content-Length: 384
Content-Type: application/json; charset=utf-8
Date: Wed, 08 Mar 2017 08:43:14 GMT

[
    {
        "ChallengeID": "0e1e4767-6df9-4869-7936-5eed4db95dde", 
        "CreateTime": "2017-02-24T20:19:52.966+08:00", 
        "FinishTime": "0001-01-01T00:00:00Z", 
        "Flag": "flag{Jayden-Thompson-December-Sunday}", 
        "Services": null, 
        "State": "failed", 
        "Services": [
            {
                "PublishedPort": 0, 
                "ServiceName": "busybox", 
                "TargetPort": 0
            }, 
            {
                "PublishedPort": 30000, 
                "ServiceName": "nginx", 
                "TargetPort": 80
            }
        ], 
        "TemplateID": "5ba174a1-cb81-4227-5f65-2a6c7985f6ea"
    }
]

# type=1 : terminated challenges
➜  ~ http "127.0.0.1:8080/user/ming-zi/challenges?type=1"
HTTP/1.1 200 OK
Content-Length: 384
Content-Type: application/json; charset=utf-8
Date: Wed, 08 Mar 2017 08:43:14 GMT

[
    {
        "ChallengeID": "671e4166-db2e-416f-7e33-898ea44e07c2", 
        "CreateTime": "2017-03-01T21:51:53.361+08:00", 
        "FinishTime": "0001-01-01T00:00:00Z", 
        "Flag": "flag{Sofia-Moore-October-Monday}", 
        "Services": [
            {
                "PublishedPort": 0, 
                "ServiceName": "busybox", 
                "TargetPort": 0
            }, 
            {
                "PublishedPort": 30000, 
                "ServiceName": "nginx", 
                "TargetPort": 80
            }
        ], 
        "State": "terminated", 
        "TemplateID": "5ba174a1-cb81-4227-5f65-2a6c7985f6ea"
    }
]


# type=2 : running challenges
➜  ~ http "127.0.0.1:8080/user/ming-zi/challenges?type=2"
HTTP/1.1 200 OK
Content-Length: 384
Content-Type: application/json; charset=utf-8
Date: Wed, 08 Mar 2017 08:43:14 GMT

[
    {
        "ChallengeID": "671e4166-db2e-416f-7e33-898ea44e07c2", 
        "CreateTime": "2017-03-01T21:51:53.361+08:00", 
        "FinishTime": "0001-01-01T00:00:00Z", 
        "Flag": "flag{Sofia-Moore-October-Monday}", 
        "Services": [
            {
                "PublishedPort": 0, 
                "ServiceName": "busybox", 
                "TargetPort": 0
            }, 
            {
                "PublishedPort": 30000, 
                "ServiceName": "nginx", 
                "TargetPort": 80
            }
        ], 
        "State": "running", 
        "TemplateID": "5ba174a1-cb81-4227-5f65-2a6c7985f6ea"
    }
]
```


### `GET    /user/:userURL/followers`
```bash
➜  ~ http "127.0.0.1:8080/user/ming-zi/followers"
HTTP/1.1 200 OK
Content-Length: 17
Content-Type: application/json; charset=utf-8
Date: Wed, 08 Mar 2017 08:49:30 GMT

{
    "followers": []
}
```

### `GET    /user/:userURL/followees`
```bash
➜  ~ http "127.0.0.1:8080/user/ming-zi/followees"
HTTP/1.1 200 OK
Content-Length: 17
Content-Type: application/json; charset=utf-8
Date: Wed, 08 Mar 2017 08:49:30 GMT

{
    "followees": []
}
```

## Challlenge
### `GET    /challenges`
```bash
# get all challenges
➜  ~ http "127.0.0.1:8080/challenges"            
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Wed, 08 Mar 2017 08:57:15 GMT
Transfer-Encoding: chunked

[
    {
        "Flag": "flag{Elijah-Miller-September-Thursday}", 
        "ID": "b1a6c487-c10e-4ac9-761e-a5f6f6b4f19c", 
        "Name": "testTemplate", 
        "StackID": "b1a6c487-c10e-4ac9-761e-a5f6f6b4f19c", 
        "State": "failed", 
        "TemplateID": "5ba174a1-cb81-4227-5f65-2a6c7985f6ea", 
        "Time": {
            "CreateTime": "2017-02-24T20:18:05.64+08:00", 
            "FinishTime": "2017-02-24T20:18:05.671+08:00"
        }, 
        "UserID": "537326cd-a113-4dbf-49f9-93234ec8799a"
    }, 
    {
        "Flag": "flag{Jayden-Thompson-December-Sunday}", 
        "ID": "0e1e4767-6df9-4869-7936-5eed4db95dde", 
        "Name": "testTemplate", 
        "StackID": "0e1e4767-6df9-4869-7936-5eed4db95dde", 
        "State": "terminated", 
        "TemplateID": "5ba174a1-cb81-4227-5f65-2a6c7985f6ea", 
        "Time": {
            "CreateTime": "2017-02-24T20:19:52.966+08:00", 
            "FinishTime": "2017-02-24T20:20:12.843+08:00"
        }, 
        "UserID": "537326cd-a113-4dbf-49f9-93234ec8799a"
    }, 
    {
        "Flag": "flag{Addison-Anderson-September-Sunday}", 
        "ID": "2f8bdc4b-2802-4e6b-5e22-5e5f0e289ca1", 
        "Name": "testTemplate", 
        "StackID": "2f8bdc4b-2802-4e6b-5e22-5e5f0e289ca1", 
        "State": "terminated", 
        "TemplateID": "5ba174a1-cb81-4227-5f65-2a6c7985f6ea", 
        "Time": {
            "CreateTime": "2017-02-24T20:20:38.409+08:00", 
            "FinishTime": "2017-02-24T20:24:16.903+08:00"
        }, 
        "UserID": "537326cd-a113-4dbf-49f9-93234ec8799a"
    }, 
    {
        "Flag": "flag{Ella-Taylor-July-Wednesday}", 
        "ID": "83184fb3-1c56-4737-6510-663a310071b1", 
        "Name": "testTemplate", 
        "StackID": "83184fb3-1c56-4737-6510-663a310071b1", 
        "State": "terminated", 
        "TemplateID": "5ba174a1-cb81-4227-5f65-2a6c7985f6ea", 
        "Time": {
            "CreateTime": "2017-02-24T20:24:19.974+08:00", 
            "FinishTime": "2017-02-24T20:37:50.67+08:00"
        }, 
        "UserID": "537326cd-a113-4dbf-49f9-93234ec8799a"
    }, 
    {
        "Flag": "flag{Abigail-Wilson-September-Thursday}", 
        "ID": "9266e1aa-6c81-45e6-6ca7-afe065010a25", 
        "Name": "testTemplate", 
        "StackID": "9266e1aa-6c81-45e6-6ca7-afe065010a25", 
        "State": "terminated", 
        "TemplateID": "5ba174a1-cb81-4227-5f65-2a6c7985f6ea", 
        "Time": {
            "CreateTime": "2017-02-24T20:38:07.121+08:00", 
            "FinishTime": "2017-02-24T20:39:01.517+08:00"
        }, 
        "UserID": "537326cd-a113-4dbf-49f9-93234ec8799a"
    }, 
    {
        "Flag": "flag{Andrew-Thomas-November-Tuesday}", 
        "ID": "12ca291f-48cf-4df0-7868-682552100e08", 
        "Name": "testTemplate", 
        "StackID": "12ca291f-48cf-4df0-7868-682552100e08", 
        "State": "terminated", 
        "TemplateID": "5ba174a1-cb81-4227-5f65-2a6c7985f6ea", 
        "Time": {
            "CreateTime": "2017-02-24T20:39:07.006+08:00", 
            "FinishTime": "2017-03-01T21:13:28.73+08:00"
        }, 
        "UserID": "537326cd-a113-4dbf-49f9-93234ec8799a"
    }, 
    {
        "Flag": "flag{Joshua-Jackson-April-Friday}", 
        "ID": "866525e0-69ed-4bdd-65d3-8466e9308eca", 
        "Name": "testTemplate", 
        "StackID": "866525e0-69ed-4bdd-65d3-8466e9308eca", 
        "State": "terminated", 
        "TemplateID": "5ba174a1-cb81-4227-5f65-2a6c7985f6ea", 
        "Time": {
            "CreateTime": "2017-03-01T21:13:45.505+08:00", 
            "FinishTime": "2017-03-01T21:18:02.792+08:00"
        }, 
        "UserID": "537326cd-a113-4dbf-49f9-93234ec8799a"
    }, 
    {
        "Flag": "flag{Daniel-Smith-July-Monday}", 
        "ID": "aed54d5e-8e77-4810-4e8e-e078a37e894c", 
        "Name": "testTemplate", 
        "StackID": "aed54d5e-8e77-4810-4e8e-e078a37e894c", 
        "State": "terminated", 
        "TemplateID": "5ba174a1-cb81-4227-5f65-2a6c7985f6ea", 
        "Time": {
            "CreateTime": "2017-03-01T21:18:39.638+08:00", 
            "FinishTime": "2017-03-01T21:27:20.442+08:00"
        }, 
        "UserID": "537326cd-a113-4dbf-49f9-93234ec8799a"
    }, 
    {
        "Flag": "flag{Lily-Johnson-November-Saturday}", 
        "ID": "a108d626-3404-4590-7e69-12a7383f60c4", 
        "Name": "testTemplate", 
        "StackID": "a108d626-3404-4590-7e69-12a7383f60c4", 
        "State": "terminated", 
        "TemplateID": "5ba174a1-cb81-4227-5f65-2a6c7985f6ea", 
        "Time": {
            "CreateTime": "2017-03-01T21:30:17.66+08:00", 
            "FinishTime": "2017-03-01T21:32:13.528+08:00"
        }, 
        "UserID": "537326cd-a113-4dbf-49f9-93234ec8799a"
    }, 
    {
        "Flag": "flag{Chloe-Williams-March-Saturday}", 
        "ID": "0d905713-a1af-4708-4911-4366aa49f8b1", 
        "Name": "testTemplate", 
        "StackID": "0d905713-a1af-4708-4911-4366aa49f8b1", 
        "State": "terminated", 
        "TemplateID": "5ba174a1-cb81-4227-5f65-2a6c7985f6ea", 
        "Time": {
            "CreateTime": "2017-03-01T21:40:16.392+08:00", 
            "FinishTime": "2017-03-01T21:42:16.267+08:00"
        }, 
        "UserID": "537326cd-a113-4dbf-49f9-93234ec8799a"
    }, 
    {
        "Flag": "flag{Noah-Taylor-August-Sunday}", 
        "ID": "c5ffe088-fee8-43a8-6e0e-425e35d230af", 
        "Name": "testTemplate", 
        "StackID": "c5ffe088-fee8-43a8-6e0e-425e35d230af", 
        "State": "terminated", 
        "TemplateID": "5ba174a1-cb81-4227-5f65-2a6c7985f6ea", 
        "Time": {
            "CreateTime": "2017-03-01T21:47:17.633+08:00", 
            "FinishTime": "2017-03-01T21:48:17.803+08:00"
        }, 
        "UserID": "537326cd-a113-4dbf-49f9-93234ec8799a"
    }, 
    {
        "Flag": "flag{Sofia-Moore-October-Monday}", 
        "ID": "671e4166-db2e-416f-7e33-898ea44e07c2", 
        "Name": "testTemplate", 
        "StackID": "671e4166-db2e-416f-7e33-898ea44e07c2", 
        "State": "running", 
        "TemplateID": "5ba174a1-cb81-4227-5f65-2a6c7985f6ea", 
        "Time": {
            "CreateTime": "2017-03-01T21:51:53.361+08:00", 
            "FinishTime": "2017-03-01T21:52:17.573+08:00"
        }, 
        "UserID": "537326cd-a113-4dbf-49f9-93234ec8799a"
    }
]
```

### `GET    /challenge/:challengeID`
```bash
# get a single challenge
➜  ~ http "127.0.0.1:8080/challenge/b1a6c487-c10e-4ac9-761e-a5f6f6b4f19c" 
HTTP/1.1 200 OK
Content-Length: 380
Content-Type: application/json; charset=utf-8
Date: Wed, 08 Mar 2017 08:58:05 GMT

{
    "Flag": "flag{Elijah-Miller-September-Thursday}", 
    "ID": "b1a6c487-c10e-4ac9-761e-a5f6f6b4f19c", 
    "Name": "testTemplate", 
    "StackID": "b1a6c487-c10e-4ac9-761e-a5f6f6b4f19c", 
    "State": "failed", 
    "TemplateID": "5ba174a1-cb81-4227-5f65-2a6c7985f6ea", 
    "Time": {
        "CreateTime": "2017-02-24T20:18:05.64+08:00", 
        "FinishTime": "2017-02-24T20:18:05.671+08:00"
    }, 
    "UserID": "537326cd-a113-4dbf-49f9-93234ec8799a"
}
```

### `POST   /challenge/create`
```bash
curl -X POST -F "uid=537326cd-a113-4dbf-49f9-93234ec8799a" -F "templateID=5ba174a1-cb81-4227-5f65-2a6c7985f6ea" "http://127.0.0.1:8080/challenge/create"

{
  "challenge created": "ok",
  "id": "2edff49e-8d21-456f-73f1-fdf4cdba9407"
}
```

### `DELETE /challenge/remove`
```bash
curl -X DELETE -F "uid=537326cd-a113-4dbf-49f9-93234ec8799a" -F "cid=2edff49e-8d21-456f-73f1-fdf4cdba9407" "http://127.0.0.1:8080/challenge/remove"

{
  "remove challenge": "ok"
}
```


## Templates
### `GET    /templates`
```bash
➜  ~ http 127.0.0.1:8080/templates                                                                            
HTTP/1.1 200 OK
Content-Length: 220
Content-Type: application/json; charset=utf-8
Date: Wed, 08 Mar 2017 09:11:52 GMT

[
    {
        "Content": null, 
        "ID": "fff", 
        "Name": "testTemplatettt"
    }, 
    {
        "Content": null, 
        "ID": "5ba174a1-cb81-4227-5f65-2a6c7985f6ea", 
        "Name": "testTemplate"
    }, 
    {
        "Content": null, 
        "ID": "ae70eeed-68b4-42ac-4ca2-17ed74289272", 
        "Name": "nginx测试"
    }
]

```

### `GET    /template/:templateID`
```bash
➜  ~ http 127.0.0.1:8080/template/5ba174a1-cb81-4227-5f65-2a6c7985f6ea
HTTP/1.1 200 OK
Content-Length: 315
Content-Type: application/json; charset=utf-8
Date: Wed, 08 Mar 2017 09:09:35 GMT

{
    "Content": "version: \"3\"\nservices:\n  nginx:\n    image: nginx\n    environment:\n      - FLAG=<FLAG>\n    ports:\n      - 80\n  busybox:\n    image: busybox\n    environment:\n      - FLAG=<FLAG>\n    command: sleep 999999", 
    "ID": "5ba174a1-cb81-4227-5f65-2a6c7985f6ea", 
    "Name": "testTemplate"
}

```

### `POST   /template/create`
```bash
➜  ~ curl 127.0.0.1:8080/template/create -X POST --form "name=nginx测试" -F "upload=@test.yml" 
{"InsertTemplate":"ok"}
```

### `DELETE /template/remove`
```
➜  ~ curl 127.0.0.1:8080/template/remove -X DELETE --form "id=ae70eeed-68b4-42ac-4ca2-17ed74289272"      
{"Remove Template":"ok"}
➜  ~ 

➜  ~ curl 127.0.0.1:8080/template/remove -X DELETE --form "id=ae70eeed-68b4-42ac-4ca2-17ed74289272" 
{"err":"not found"}
```

## Service
### `GET    /services`
```bash
➜  ~ http 127.0.0.1:8080/services                                                                        
HTTP/1.1 200 OK
Content-Length: 94
Content-Type: application/json; charset=utf-8
Date: Wed, 08 Mar 2017 09:13:27 GMT

[
    "671e4166-db2e-416f-7e33-898ea44e07c2_nginx", 
    "671e4166-db2e-416f-7e33-898ea44e07c2_busybox"
]

```

### `GET    /service/:serviceID`
```bash
➜  ~ http 127.0.0.1:8080/service/671e4166-db2e-416f-7e33-898ea44e07c2_nginx
HTTP/1.1 200 OK
Content-Length: 1229
Content-Type: application/json; charset=utf-8
Date: Wed, 08 Mar 2017 09:15:39 GMT

{
    "CreatedAt": "2017-03-01T13:51:56.545888793Z", 
    "Endpoint": {
        "Ports": [
            {
                "Protocol": "tcp", 
                "PublishMode": "ingress", 
                "PublishedPort": 30000, 
                "TargetPort": 80
            }
        ], 
        "Spec": {
            "Mode": "vip", 
            "Ports": [
                {
                    "Protocol": "tcp", 
                    "PublishMode": "ingress", 
                    "TargetPort": 80
                }
            ]
        }, 
        "VirtualIPs": [
            {
                "Addr": "10.255.0.2/16", 
                "NetworkID": "onutarvaid4yerz3ff59etdli"
            }, 
            {
                "Addr": "10.0.0.2/24", 
                "NetworkID": "i8be65ajb0dzo2ou4wlem59ec"
            }
        ]
    }, 
    "ID": "c606sz10sxwthkcm2ha6r7l1z", 
    "Spec": {
        "EndpointSpec": {
            "Mode": "vip", 
            "Ports": [
                {
                    "Protocol": "tcp", 
                    "PublishMode": "ingress", 
                    "TargetPort": 80
                }
            ]
        }, 
        "Labels": {
            "com.docker.stack.namespace": "671e4166-db2e-416f-7e33-898ea44e07c2"
        }, 
        "Mode": {
            "Replicated": {
                "Replicas": 1
            }
        }, 
        "Name": "671e4166-db2e-416f-7e33-898ea44e07c2_nginx", 
        "Networks": [
            {
                "Aliases": [
                    "nginx"
                ], 
                "Target": "i8be65ajb0dzo2ou4wlem59ec"
            }
        ], 
        "TaskTemplate": {
            "ContainerSpec": {
                "Env": [
                    "FLAG=flag{Sofia-Moore-October-Monday}"
                ], 
                "Image": "nginx:latest@sha256:4296639ebdf92f035abf95fee1330449e65990223c899838283c9844b1aaac4c", 
                "Labels": {
                    "com.docker.stack.namespace": "671e4166-db2e-416f-7e33-898ea44e07c2"
                }
            }, 
            "ForceUpdate": 0, 
            "Placement": {}, 
            "Resources": {}
        }
    }, 
    "UpdateStatus": {
        "CompletedAt": "0001-01-01T00:00:00Z", 
        "StartedAt": "0001-01-01T00:00:00Z"
    }, 
    "UpdatedAt": "2017-03-03T03:30:52.886026574Z", 
    "Version": {
        "Index": 14249
    }
}



```

### `GET    /service/:serviceID/status`
```bash
➜  ~ http 127.0.0.1:8080/service/671e4166-db2e-416f-7e33-898ea44e07c2_nginx/status
HTTP/1.1 200 OK
Content-Length: 1995
Content-Type: application/json; charset=utf-8
Date: Wed, 08 Mar 2017 09:14:40 GMT

{
    "CreatedAt": "2017-03-03T03:31:02.503193615Z", 
    "DesiredState": "running", 
    "ID": "6v3ofpzodqepkas01s3e7qkvx", 
    "NetworksAttachments": [
        {
            "Addresses": [
                "10.255.0.4/16"
            ], 
            "Network": {
                "CreatedAt": "2016-12-11T08:29:09.618946457Z", 
                "DriverState": {
                    "Name": "overlay", 
                    "Options": {
                        "com.docker.network.driver.overlay.vxlanid_list": "4096"
                    }
                }, 
                "ID": "onutarvaid4yerz3ff59etdli", 
                "IPAMOptions": {
                    "Configs": [
                        {
                            "Gateway": "10.255.0.1", 
                            "Subnet": "10.255.0.0/16"
                        }
                    ], 
                    "Driver": {
                        "Name": "default"
                    }
                }, 
                "Spec": {
                    "DriverConfiguration": {}, 
                    "IPAMOptions": {
                        "Configs": [
                            {
                                "Gateway": "10.255.0.1", 
                                "Subnet": "10.255.0.0/16"
                            }
                        ], 
                        "Driver": {}
                    }, 
                    "Labels": {
                        "com.docker.swarm.internal": "true"
                    }, 
                    "Name": "ingress"
                }, 
                "UpdatedAt": "2017-03-03T03:30:52.882798873Z", 
                "Version": {
                    "Index": 14246
                }
            }
        }, 
        {
            "Addresses": [
                "10.0.0.5/24"
            ], 
            "Network": {
                "CreatedAt": "2017-03-01T13:51:53.487289337Z", 
                "DriverState": {
                    "Name": "overlay", 
                    "Options": {
                        "com.docker.network.driver.overlay.vxlanid_list": "4097"
                    }
                }, 
                "ID": "i8be65ajb0dzo2ou4wlem59ec", 
                "IPAMOptions": {
                    "Configs": [
                        {
                            "Gateway": "10.0.0.1", 
                            "Subnet": "10.0.0.0/24"
                        }
                    ], 
                    "Driver": {
                        "Name": "default"
                    }
                }, 
                "Spec": {
                    "DriverConfiguration": {
                        "Name": "overlay"
                    }, 
                    "IPAMOptions": {
                        "Driver": {}
                    }, 
                    "Labels": {
                        "com.docker.stack.namespace": "671e4166-db2e-416f-7e33-898ea44e07c2"
                    }, 
                    "Name": "671e4166-db2e-416f-7e33-898ea44e07c2_default"
                }, 
                "UpdatedAt": "2017-03-03T03:30:52.884029592Z", 
                "Version": {
                    "Index": 14247
                }
            }
        }
    ], 
    "NodeID": "n6mwx010i28dtjrggiy7zfk11", 
    "ServiceID": "c606sz10sxwthkcm2ha6r7l1z", 
    "Slot": 1, 
    "Spec": {
        "ContainerSpec": {
            "Env": [
                "FLAG=flag{Sofia-Moore-October-Monday}"
            ], 
            "Image": "nginx:latest@sha256:4296639ebdf92f035abf95fee1330449e65990223c899838283c9844b1aaac4c", 
            "Labels": {
                "com.docker.stack.namespace": "671e4166-db2e-416f-7e33-898ea44e07c2"
            }
        }, 
        "ForceUpdate": 0, 
        "Placement": {}, 
        "Resources": {}
    }, 
    "Status": {
        "ContainerStatus": {
            "ContainerID": "4b4cfe6a026b9d437ab34314c66cad70265960a94d37d25b94761ab1cc5dd357", 
            "PID": 22423
        }, 
        "Message": "started", 
        "PortStatus": {}, 
        "State": "running", 
        "Timestamp": "2017-03-03T03:31:09.250673241Z"
    }, 
    "UpdatedAt": "2017-03-03T03:31:09.33761384Z", 
    "Version": {
        "Index": 14263
    }
}

```