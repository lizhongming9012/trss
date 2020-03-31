package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type SysInfo struct {
	Uri      string `yaml:"uri"`
	Username string `yaml:"username"`
	Pwd      string `yaml:"pwd"`
}
type DbInfo struct {
	Uri      string `yaml:"uri"`
	User     string `yaml:"user"`
	Pwd      string `yaml:"pwd"`
	Port     string `yaml:"port"`
	DataName string `yaml:"dataname"`
}

type Conf struct {
	Sysinfo      SysInfo `yaml:"sysinfo"`
	Dbinfo       DbInfo  `yaml:"dbinfo"`
	GoroutineNum int     `yaml:"goroutinenum"`
	Limit        int     `yaml:"limit"`
	LogFile      string  `yaml:"logfile"`
}

func (c *Conf) ConfReader() {
	//c := C{}
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	err = yaml.Unmarshal([]byte(yamlFile), &c)
	if err != nil {
		log.Fatalf("error:%v", err)
	}
	fmt.Printf("sys uri is :%s,username is :%s,pwd is :%s\n", c.Sysinfo.Uri, c.Sysinfo.Username, c.Sysinfo.Pwd)
	fmt.Printf("db uri is :%s,user is :%s,pwd is :%s,port is :%s,databasename is %s\n", c.Dbinfo.Uri, c.Dbinfo.User, c.Dbinfo.Pwd, c.Dbinfo.Port, c.Dbinfo.DataName)

}
