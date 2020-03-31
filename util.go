package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net"
	"strings"
	"time"
)

// sleep n seconds
func Delay(n int) {
	time.Sleep(time.Duration(n) * time.Second)
}

//iteration the p and check if value are key in map,otherwise fill the key with blank string
func Preserve(p []string, params map[string]string) {
	for _, s := range p {
		_, ok := params[s]
		if !ok {
			params[s] = ""
		}
	}
}

// return the local host' ip and mac address
func IpAndMac() (string, string, error) {
	addrs, err := net.InterfaceAddrs()
	var ip, mac string
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	var currentIP, currentNetworkHardwareName string

	for _, address := range addrs {

		// check the address type and if it is not a loopback the display it
		// = GET LOCAL IP ADDRESS
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println("Current IP address : ", ipnet.IP.String())
				currentIP = ipnet.IP.String()
				ip = currentIP
			}
		}
	}
	// get all the system's or local machine's network interfaces
	interfaces, _ := net.Interfaces()
	for _, interf := range interfaces {

		if addrs, err := interf.Addrs(); err == nil {
			for index, addr := range addrs {
				fmt.Println("[", index, "]", interf.Name, ">", addr)

				// only interested in the name with current IP address
				if strings.Contains(addr.String(), currentIP) {
					fmt.Println("Use name : ", interf.Name)
					currentNetworkHardwareName = interf.Name
				}
			}
		}
	}
	// extract the hardware information base on the interface name
	// capture above
	netInterface, err := net.InterfaceByName(currentNetworkHardwareName)

	if err != nil {
		fmt.Println(err)
		return ip, "", err
	}
	macAddress := netInterface.HardwareAddr
	// verify if the MAC address can be parsed properly
	hwAddr, err := net.ParseMAC(macAddress.String())

	if err != nil {
		fmt.Println("No able to parse MAC address : ", err)
		return ip, "", err
	}
	mac = hwAddr.String()
	return ip, mac, nil
}
func getNodeNameValue(s *goquery.Selection) (string, string) {
	var k, vv string
	for _, v := range s.Nodes[0].Attr {
		if strings.Compare(v.Key, "name") == 0 {
			k = v.Val
		}
		if strings.Compare(v.Key, "value") == 0 {
			vv = v.Val
		}
	}
	return k, vv
}
func getNodeAttrVal(s *goquery.Selection, attr string) string {
	var vv string
	for _, v := range s.Nodes[0].Attr {
		if strings.Compare(v.Key, attr) == 0 {
			vv = v.Val
			break
		}
	}
	return vv
}
func CheckErrs(errs []error, uri string) {
	if len(errs) != 0 {
		for _, e := range errs {
			log.Println(e.Error())
		}
		log.Fatalf("request errors.for the request:%s", uri)
	}
}
