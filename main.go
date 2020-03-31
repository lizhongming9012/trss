package main

import (
	"github.com/parnurzeal/gorequest"
	"io"
	"log"
	"os"
)

type StringMap map[string]string
type StringList []map[string]string

var swrydm string
var c = Conf{}
var local = "1370683"
var ip, mac string

var IeHeader = "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0; .NET4.0C; .NET4.0E; .NET CLR 2.0.50727; .NET CLR 3.0.4506.2152; .NET CLR 3.5.30729; InfoPath.2)"

func main() {
	var logWriter, err = os.Create(c.LogFile)
	if err != nil {
		log.Println("create log file err...")
		os.Exit(4)
	}
	var mw = io.MultiWriter(os.Stdout, logWriter)
	log.SetOutput(mw)
	ip, mac, err = IpAndMac()
	if err != nil {
		log.Fatalf("get ip and mac fail...")
		os.Exit(2)
	}
	var client = gorequest.New()
	if login(client) {
		swdjService(client)
	}
	log.Println("----that's ok!---")
}

func init() {

	c.ConfReader()
}
