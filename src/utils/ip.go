package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
)

func TaobaoIP2Region(ip string) (string, error) {
	Region := ""
	URL := "http://ip.taobao.com/service/getIpInfo.php?ip=" + ip
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Get(URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	jsBody, _ := simplejson.NewJson([]byte(body))
	code, _ := jsBody.Get("code").Int()
	switch code {
	case 0:
		country, _ := jsBody.Get("data").Get("country").String()
		region, _ := jsBody.Get("data").Get("region").String()
		city, _ := jsBody.Get("data").Get("city").String()
		isp, _ := jsBody.Get("data").Get("isp").String()
		if isp != "" {
			Region = country + region + city + "[" + isp + "]"
		} else {
			Region = country + region + city
		}
	case 1:
		Region, _ = jsBody.Get("data").String()
	}
	return Region, nil
}

func BaiduIP2Region(ip string) (string, error) {
	// 'https://sp1.baidu.com/8aQDcjqpAAV3otqbppnN2DJv/api.php?query=60.117.48.212&resource_id=6006&ie=utf8&oe=utf8&format=json'
	var err error
	region := &simplejson.Json{}
	protocal := "https://"
	hosts := []string{"sp0", "sp1", "sp2", "sp3"}
	url := ".baidu.com/8aQDcjqpAAV3otqbppnN2DJv/api.php?query=" + ip + "&resource_id=6006&ie=utf8&oe=utf8&format=json"
	for _, host := range hosts {
		URL := protocal + host + url
		client := &http.Client{
			Timeout: 3 * time.Second,
		}
		resp, err := client.Get(URL)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		jsBody, _ := simplejson.NewJson([]byte(body))
		jsData := jsBody.Get("data").GetIndex(0)
		region = jsData.Get("location")
		break
	}
	if err != nil {
		return "", err
	}
	r, _ := region.String()
	return r, nil
}

func OpenGPSIP2Region(ip string) (string, error) {
	Region := ""
	URL := "https://www.opengps.cn/Data/IP/IPLocHiAcc.ashx"
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	postBody := strings.NewReader("ip=" + ip)
	bodyType := "application/x-www-form-urlencoded"
	resp, err := client.Post(URL, bodyType, postBody)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	jsBody, _ := simplejson.NewJson([]byte(body))
	success, _ := jsBody.Get("success").Bool()
	if !success {
		Region = "err"
	}
	values := jsBody.Get("values").GetIndex(0)
	resultCode, _ := values.Get("resultCode").String()
	if resultCode != "定位成功" {
		return "", errors.New(resultCode)
	}
	Region, err = values.Get("address").String()
	return Region, err
}

func IP2Region(ip string) string {
	ch := make(chan string, 3)
	defer close(ch)

	go func(ch chan string) {
		r, err := OpenGPSIP2Region(ip)
		if err != nil {
			ch <- ""
			return
		}
		ch <- r
	}(ch)

	go func(ch chan string) {
		r, err := BaiduIP2Region(ip)
		if err != nil {
			ch <- ""
			return
		}
		ch <- r
	}(ch)

	go func(ch chan string) {
		r, err := TaobaoIP2Region(ip)
		if err != nil {
			ch <- ""
			return
		}
		ch <- r
	}(ch)

	for {
		if len(ch) == 3 {
			break
		}
		time.Sleep(100 * time.Microsecond)
	}

	return fmt.Sprintf("[%s];[%s];[%s]", <-ch, <-ch, <-ch)
}
