package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
	"log"
	"strings"
)

var n = 1
var ck string

func login(client *gorequest.SuperAgent) bool {

	//1.get request to open login page
	params := firstReqGet(client)
	if len(params) == 0 {
		log.Fatalf("step 1:query params err")
	}
	Delay(n)
	//2.post request to post login submit
	pa := secondReqPost(client, params)
	if len(pa) == 0 {
		log.Fatalf("step 2:seek input err")
	}
	Delay(n)
	//3.post js_sso_server chooseIdentify
	thirdReqPost(client, pa)
	Delay(n)
	//4.get request swrydm
	swrydm = forthReqGet(client)
	if len(swrydm) == 11 {
		fmt.Println("登陆金税三期核心征管成功!")
		return true
	}
	return false
}
func firstReqGet(client *gorequest.SuperAgent) StringMap {
	params := make(StringMap)
	preserve := []string{"brower", "hostname", "app", "smscode", "phonenum", "yzm", "swrydm", "lt", "execution", "_eventId", "submit"}
	resp, b, errs := client.Get(c.Sysinfo.Uri).Set("User-Agent", IeHeader).End()
	CheckErrs(errs, c.Sysinfo.Uri)
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Parse the body fail...%v", err)
		return params
	}
	doc.Find("#_user_div_login_").Find("input").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the name and value
		k, vv := getNodeNameValue(s)
		params[k] = vv
	})
	ck = resp.Header.Get("Set-Cookie")
	if strings.Index(ck, ";") <= 1 {
		log.Println(resp.Header)
		log.Fatalln(b)
	}
	ck = ck[:strings.Index(ck, ";")-1]
	params["username"] = c.Sysinfo.Username
	params["password"] = mdf5(c.Sysinfo.Pwd)
	if err != nil {
		log.Fatalf("get address and mac error .%v", err)
		return params
	}
	params["ip"] = ip
	params["mac"] = mac
	params["hdsn"] = "-1962243318"
	// iteration preserve param
	Preserve(preserve, params)
	return params
}
func secondReqPost(client *gorequest.SuperAgent, params StringMap) StringMap {
	uri := "http://dddl.sdsw.tax.cn/js_sso_server/login?service=http%3A%2F%2Fportal.sdsw.tax.cn%2Fwelcome.html"
	resp, b, errs := client.Post(uri).
		Set("Accept", "image/jpeg, application/x-ms-application, image/gif, application/xaml+xml, image/pjpeg, application/x-ms-xbap, application/vnd.ms-excel, application/vnd.ms-powerpoint, application/msword, */*").
		Set("Accept-Language", "zh-CN").
		Set("User-Agent", IeHeader).
		Type(gorequest.TypeForm).
		SendMap(params).
		End()
	CheckErrs(errs, uri)
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("second->Parse the body fail...%v", err)
	}
	nd := doc.Find("#msg").Nodes
	if len(nd) > 0 {
		log.Println("SECOND:", b)
		log.Fatalln("post msg is \n", nd[0].Attr)
	}
	param2s := make(StringMap)
	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the name and value
		k, vv := getNodeNameValue(s)
		param2s[k] = vv
	})
	return param2s
}
func thirdReqPost(client *gorequest.SuperAgent, params StringMap) bool {
	uri := "http://dddl.sdsw.tax.cn/js_sso_server/chooseIdentify"
	// iteration preserve param
	_, _, errs := client.Post(uri).Set("Accept", "image/jpeg, application/x-ms-application, image/gif, application/xaml+xml, image/pjpeg, application/x-ms-xbap, application/vnd.ms-excel, application/vnd.ms-powerpoint, application/msword, */*").
		Set("Accept-Language", "zh-CN").
		Set("User-Agent", IeHeader).
		Type(gorequest.TypeForm).
		SendMap(params).End()
	if len(errs) != 0 {
		log.Fatalf("chooseIdentify err")
	} else {
		return true
	}
	return false
}

func forthReqGet(client *gorequest.SuperAgent) string {
	uri := "http://portal.sdsw.tax.cn/sword?ctrl=MH003InitLoginxxCtrl_openWin"
	var s, res string
	resp, b, errs := client.Get(uri).End()
	CheckErrs(errs, uri)
	fmt.Println("five:", b)
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("five->Parse the body fail...%v", err)
		return res
	}
	attr := doc.Find("#SwordPageData").Nodes[0].Attr
	for _, v := range attr {
		if strings.Compare(v.Key, "data") == 0 {
			s = v.Val
			break
		}
	}
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(s), &dat); err != nil {
		log.Fatalf("sixth req:parse data to json err")
		return res
	}
	ss := dat["data"].([]interface{})
	sss := ss[0].(map[string]interface{})
	data := sss["data"].(map[string]interface{})
	swry := data["swrydm"].(map[string]interface{})
	result := swry["value"]
	res = result.(string)
	return res
}
