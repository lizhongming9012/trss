package main

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type SwdjPostBody struct {
	sqlxh           string
	tj              string
	fzzd            string
	znl             string
	dmtomc          string
	totalCount      string
	totalflag       string
	gwssswjg1       string
	qxswjgstr       string
	cxsql           string
	footstr         string
	footvalue       string
	fuzzyQueryTj    string
	sjly            string
	pagecount_bdxt  string
	bdxt_totalcount string
	qshs_dfxt       string
	page            string
	start           string
	limit           string
}

var paramswdj map[string]string
var insertChan chan StringList
var quit chan int

func swdjService(cli *gorequest.SuperAgent) {

	go prepare(cli)
	func() {
		var m []map[string]string
		for {
			select {
			case v := <-insertChan:
				m = append(m, v...)
				if len(m) >= 500 {
					dbInsert(m, "gt3_dj_nsrxx")
					m = []map[string]string{}
				}
			case <-quit:
				if len(m) != 0 {
					dbInsert(m, "gt3_dj_nsrxx")
					m = []map[string]string{}
				}
				return
			}
		}
	}()
}
func prepare(cli *gorequest.SuperAgent) {
	swjgArray := []string{"13706830000", "13706340000"}
	u := `http://tycx.sdsw.tax.cn/sword?swordsc=NudM%25252Fj2%2525252BJUy8nDbt25NS8g%25253D%25253D&sqlxh=10010002&rUUID=iu4twNVlWNeITdFpQoMPIvLM5PntifqU&zndm=01&gndm=s%25252525%25252523wordMkb4anpGJ16t3kbh84iWVNt44NQHTLzn&gwssswjg=s%25252525%25252523wordxYXlAvIBqCs%2525252Bwa8lIKjaMA%25253D%25253D&gwxh=3A6917965FC81188E05311610C4C9131&gt3zyydm=F0000003&gt3ywfldm=null&gbiSage=CUR`
	params := make(map[string]string)
	p := []string{"service", "swrysfdm"}
	resp, b, errs := cli.Get(u).Set("User-Agent", IeHeader).End()
	CheckErrs(errs, c.Sysinfo.Uri)
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Parse the body fail...%v", err)
	}
	inputs := doc.Find("input")
	if len(inputs.Nodes) == 0 {
		log.Println(b)
		log.Fatalln("query param err swdj first")
	}
	inputs.Each(func(i int, s *goquery.Selection) {
		k, vv := getNodeNameValue(s)
		params[k] = vv
	})
	for _, s := range p {
		_, ok := params[s]
		if !ok {
			params[s] = ""
		}
	}
	resp, b, errs = cli.Post("http://dddl.sdsw.tax.cn/js_sso_server/chooseIdentify").
		Set("Accept", "image/jpeg, application/x-ms-application, image/gif, application/xaml+xml, image/pjpeg, application/x-ms-xbap, application/vnd.ms-excel, application/vnd.ms-powerpoint, application/msword, */*").
		Set("Accept-Language", "zh-CN").
		Set("User-Agent", IeHeader).
		Type(gorequest.TypeForm).
		SendMap(params).End()
	CheckErrs(errs, "http://dddl.sdsw.tax.cn/js_sso_server/chooseIdentify")
	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Parse the body fail...%v", err)
	}
	if len(doc.Find("a").Nodes) == 0 {
		log.Println(b)
		log.Fatalln("chooseIdentify err...")
	}
	nextUri := getNodeAttrVal(doc.Find("a"), "href")
	if len(nextUri) < 3 {
		log.Fatalf("can't find next uri %s", nextUri)
	}
	resp, body, errs := cli.Get(nextUri).Set("User-Agent", IeHeader).End()
	beg := strings.Index(body, "[{")
	end := strings.Index(body, "}]}")
	res := body[beg : end+2]
	var dat []interface{}
	if err := json.Unmarshal([]byte(res), &dat); err != nil {
		log.Fatalf("Parse the body err:for %s,err:%v", nextUri, err)
	}
	params = make(map[string]string)
	for _, n := range dat {
		v := n.(map[string]string)
		params[v["name"]] = v["value"]
	}
	params["gwssswjg1"] = params["gwssswjg"]
	params["gwxhqs"] = params["gwxh"]
	params["zndm1"] = params["zndm"]
	params["gndm1"] = params["gndm"]
	p = []string{"crtj", "gndm1", "gwssswjg1", "gwxh", "gwxhqs", "sqlxh", "zndm1"}
	Preserve(p, params)
	resp, _, errs = cli.Post("http://tycx.sdsw.tax.cn/download.sword?ctrl=CX301CxcshCtrl_cxtjcsh&n="+string(time.Now().UnixNano())).
		Set("Accept", "*/*").
		Set("Accept-Language", "zh-CN").
		Set("User-Agent", IeHeader).
		Type("urlencoded").
		SendMap(params).End()
	CheckErrs(errs, "http://tycx.sdsw.tax.cn/download.sword?ctrl=CX301CxcshCtrl_cxtjcsh&n="+string(time.Now().UnixNano()))

	for _, swjg := range swjgArray {
		log.Printf("获取%s税务登记相关信息...", swjg)
		swjgQuery(cli, swjg)
	}
	quit <- 0
}
func swjgQuery(cli *gorequest.SuperAgent, swjg string) {
	_, _, errs := cli.Post("http://tycx.sdsw.tax.cn/download.sword?ctrl=CX302ZxcxCtrl_getFieldsAndColumns&n="+string(time.Now().UnixNano())).
		Send(`{"sqlxh":"10010002"}`).
		Send(`"cxjgl":"NSRSBH,NSRMC,NSRZT_DM,DJZCLX_DM,KYSLRQ,HY_DM,DJJG_DM,ZGSWJ_DM,ZGSWSKFJ_DM,SSGLY_DM,JDXZ_DM,FDDBRXM,FDDBRSFZJLX_DM,FDDBRSFZJHM,FDDBRYDDH"`).
		Send(`"cxlx":"1"`).
		Send(`"tj":"%5B%7Bname%3ASFBAHTDZSBM%2Ctype%3Astring%2Ctjzmerge%3Aundefined%2Cvalue%3A'N'%7D%2C%7Bname%3AYXBZ%2Ctype%3Astring%2Ctjzmerge%3A0%2Cvalue%3A'Y'%7D%2C%7Bname%3AZGSWJ_DM%2Ctype%3Astring%2Ctjzmerge%3A0%2Cvalue%3A'`+swjg+`'%7D%5D"`).
		Send(`"fzzd":""`).
		Send(`"wd":""`).
		Send(`"zb":""`).
		Set("Accept", "*/*").
		Set("User-Agent", IeHeader).
		End()
	CheckErrs(errs, "http://tycx.sdsw.tax.cn/download.sword?ctrl=CX302ZxcxCtrl_getFieldsAndColumns&n=")
	log.Printf("查询税务机关代码为%s的登记信息", swjg)
	var pages int
	var pageSize = c.Limit
	footvalue := `{"TOTALCOUNT":"+totalCount+","ZCZB":"381,692,500,655.02","TZZE":"99,740,045,086.01","GDGRS":"308778","CYRS":"562318"}`
	tot := getPageQueryPost(cli, 1, 0, swjg, "", "null")

	if tot%pageSize == 0 {
		pages = tot / pageSize
	} else {
		pages = tot/pageSize + 1
	}
	var wg sync.WaitGroup
	var seg int
	if pages%10 == 0 {
		seg = pages / 10
	} else {
		seg = pages/10 + 1
	}
	for i := 1; i <= 10; i++ {

		wg.Add(1)
		go segPageQueryPost((i-1)*seg, i*seg, pages, cli.Clone(), &wg, swjg, string(tot), footvalue)
	}
	wg.Wait()
	log.Printf("%s税务登记信息获取结束", swjg)
}
func segPageQueryPost(beg int, end int, pages int, cli *gorequest.SuperAgent, wg *sync.WaitGroup, swjg string, totStr string, footvalue string) {
	if wg != nil {
		defer wg.Done()
	}
	for j := beg; j < end; j++ {
		if j <= 1 {
			continue
		}
		if j <= pages {
			getPageQueryPost(cli.Clone(), j, (j-1)*c.Limit, swjg, totStr, footvalue)
		}
	}
}
func getPageQueryPost(cli *gorequest.SuperAgent, page int, start int, swjg string, totStr string, footvalue string) int {
	uri := `http://tycx.sdsw.tax.cn/download.sword?ctrl=CX302ZxcxCtrl_exequery&sjymc=tycx_sdswhxcx_gl&limtTime=&gwxhqs=3A6917965FC81188E05311610C4C9131&n=1553132276251&cxlx=1&wd=undefined&zb=undefined&queryplugin=&_dc=1553132276259`
	b := SwdjPostBody{sqlxh: "10010002", tj: `[{name:SFBAHTDZSBM,type:string,tjzmerge:undefined,value:'N'},{name:YXBZ,type:string,tjzmerge:0,value:'Y'},{name:ZGSWJ_DM,type:string,tjzmerge:0,value:'` + swjg + `'}]`,
		fzzd: "", znl: "null", dmtomc: `[{"dm":"DJZCLX_DM","dmweb":"DJZCLX_DM","mc":"DJZCLXMC","dmb":"DM_DJ_DJZCLX"},{"dm":"GDGHLX_DM","dmweb":"GDGHLX_DM","mc":"GDGHLXMC","dmb":"DM_DJ_GDGHLX"},{"dm":"KZZTDJLX_DM","dmweb":"KZZTDJLX_DM","mc":"KZZTDJLXMC","dmb":"DM_DJ_KZZTDJLX"},{"dm":"NSRZTLX_DM","dmweb":"NSRZTLX_DM","mc":"NSRZTLXMC","dmb":"DM_DJ_NSRZTLX"},{"dm":"SWDJBZFS_DM","dmweb":"BZFS_DM","mc":"SWDJBZFSMC","dmb":"DM_DJ_SWDJBZFS"},{"dm":"ZFJGLX_DM","dmweb":"ZFJGLX_DM","mc":"ZFJGLXMC","dmb":"DM_DJ_ZFJGLX"},{"dm":"DWLSGX_DM","dmweb":"DWLSGX_DM","mc":"DWLSGXMC","dmb":"DM_GY_DWLSGX"},{"dm":"GSXZGLJG_DM","dmweb":"PZSLJG_DM","mc":"GSXZGLJGMC","dmb":"DM_GY_GSXZGLJG"},{"dm":"GYKGLX_DM","dmweb":"GYKGLX_DM","mc":"GYKGLXMC","dmb":"DM_GY_GYKGLX"},{"dm":"HSFS_DM","dmweb":"HSFS_DM","mc":"HSFSMC","dmb":"DM_GY_HSFS"},{"dm":"HY_DM","dmweb":"HY_DM","mc":"HYMC","dmb":"DM_GY_HY"},{"dm":"JDXZ_DM","dmweb":"JDXZ_DM","mc":"JDXZMC","dmb":"DM_GY_JDXZ"},{"dm":"KJZDZZ_DM","dmweb":"KJZDZZ_DM","mc":"KJZDZZMC","dmb":"DM_GY_KJZDZZ"},{"dm":"NSRZT_DM","dmweb":"NSRZT_DM","mc":"NSRZTMC","dmb":"DM_GY_NSRZT"},{"dm":"SFBZ_DM","dmweb":"FJMQYBZ@KQCCSZTDJBZ@YXBZ@MYQYBZ","mc":"SFBZMC","dmb":"DM_GY_SFBZ"},{"dm":"SFZJLX_DM","dmweb":"FDDBRSFZJLX_DM","mc":"SFZJLXMC","dmb":"DM_GY_SFZJLX"},{"dm":"SWJG_DM","dmweb":"DJJG_DM@ZGSWJ_DM@ZGSWSKFJ_DM@PGJG_DM","mc":"SWJGMC","dmb":"DM_GY_SWJG"},{"dm":"SWRY_DM","dmweb":"SSGLY_DM@LRR_DM@XGR_DM","mc":"SWRYMC","dmb":"DM_GY_SWRY"},{"dm":"XZQHSZ_DM","dmweb":"ZCDZXZQHSZ_DM@SCJYDZXZQHSZ_DM","mc":"XZQHMC","dmb":"DM_GY_XZQH"},{"dm":"YGZNSRLX_DM","dmweb":"YGZNSRLX_DM","mc":"YGZNSRLXMC","dmb":"DM_GY_YGZNSRLX"},{"dm":"ZZJGLX_DM","dmweb":"ZZJGLX_DM","mc":"ZZJGLXMC","dmb":"DM_GY_ZZJGLX"},{"dm":"ZZLX_DM","dmweb":"ZZLX_DM","mc":"ZZLXMC","dmb":"DM_GY_ZZLX"}]`,
		totalCount: totStr, totalflag: "", gwssswjg1: "13700001800", qxswjgstr: "", cxsql: "", footstr: `to_char(SUM(TZZE) OVER(),'FM999,999,999,999,999,990.00') TZZE_ODS,to_char(SUM(ZCZB) OVER(),'FM999,999,999,999,999,990.00') ZCZB_ODS,to_char(SUM(GDGRS) OVER()) GDGRS_ODS,to_char(SUM(CYRS) OVER()) CYRS_ODS,`,
		footvalue:    footvalue,
		fuzzyQueryTj: "", sjly: "", pagecount_bdxt: "", bdxt_totalcount: "", qshs_dfxt: "",
		page: string(page), start: string(start),
		limit: string(c.Limit),
	}
	_, body, errs := cli.Post(uri).Send(b).Set("Accept", "*/*").Set("User-Agent", IeHeader).Set("x-requested-with", "XMLHttpRequest").Type("urlencoded").End()
	CheckErrs(errs, uri)
	beg := strings.Index(body, "[{")
	end := strings.Index(body, "}")
	s := body[beg:end+1] + "]"
	var dat map[string]string
	json.Unmarshal([]byte(s), &dat)
	stot := dat["totalCount"]
	tot, err := strconv.Atoi(stot)
	if !(err == nil && tot > 0) {
		log.Fatalf("query tot err.%s", uri)
	}
	topics := dat["topics"]
	var topicList []map[string]string
	json.Unmarshal([]byte(topics), &topicList)
	sendSwjgChan(topicList, swjg)
	return tot
}
func sendSwjgChan(list StringList, swjg string) {
	fields := getFields("gt3_dj_nsrxx")
	var dataset = make(StringList, len(list))
	for i, v := range list {
		for _, f := range fields {
			dataset[i][f.name] = v[f.name]
		}
		dataset[i]["swjg_dm"] = swjg
	}
	insertChan <- dataset
}
