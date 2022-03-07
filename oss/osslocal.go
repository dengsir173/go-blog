package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func main() {
	bucketname := "geleigo"
	endpoint := "oss-cn-shenzhen.aliyuncs.com"
	objectname := "11"                                  //存储路径
	accesskey := "xxx"             //您的Accesskey
	accesskeysecret := "xxx" //您的Accesskeysecret
	contenttype := "application/json"
	gmtdate := time.Now().UTC().Format(http.TimeFormat)
	stringtosgin := "PUT\n\n" + contenttype + "\n" + gmtdate + "\n" + "/" + bucketname + "/" + objectname
	// HMACSHA1 实现部分
	key := []byte(accesskeysecret)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(stringtosgin))
	//进行base64编码
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	url := "http://" + bucketname + "." + endpoint + "/" + objectname
	payload := strings.NewReader("{go:test}")
	req, _ := http.NewRequest("PUT", url, payload)
	req.Header.Add("Content-Type", contenttype)
	req.Header.Add("Authorization", "OSS "+accesskey+":"+signature)
	req.Header.Add("Date", gmtdate)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(res)
	fmt.Println(string(body))
}
