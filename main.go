package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strings"
)

type Server struct {
	name      string
	ipAddress string
}

type Snip struct {
	ipAddress  string
	subnetMask string
}

func GetFile(fileName string) (string, error) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(file), nil
}

func GetConfig(file, pattern string) ([]string, error) {
	regexer, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	results := regexer.FindAllString(file, -1)
	return results, nil
}

func RemoveConfigKeywords(textLine, pattern string) string {
	result := strings.Replace(textLine, pattern, "", 1)
	return result
}

func GetServers(fileName string) ([]Server, error) {
	var servers []Server
	file, err := GetFile(fileName)
	if err != nil {
		return nil, err
	}
	addServerLines, err := GetConfig(file, "(add server).*")
	if err != nil {
		return nil, err
	}
	for _, addServerLine := range addServerLines {
		serverLine := RemoveConfigKeywords(addServerLine, "add server ")
		serverLineArray := strings.Split(serverLine, " ")
		var server Server
		server.name = serverLineArray[0]
		server.ipAddress = strings.Replace(serverLineArray[1], "\r", "", -1)
		servers = append(servers, server)
	}
	return servers, nil
}

func GetSnips(fileName string) ([]Snip, error) {
	var snips []Snip
	file, err := GetFile(fileName)
	if err != nil {
		return nil, err
	}
	addNsIpLines, err := GetConfig(file, "(add ns ip ).*")
	if err != nil {
		return nil, err
	}
	for _, addNsIpLine := range addNsIpLines {
		nsIpLine := RemoveConfigKeywords(addNsIpLine, "add ns ip ")
		nsIpLineArray := strings.Split(nsIpLine, " ")
		var snip Snip
		snip.ipAddress = nsIpLineArray[0]
		snip.subnetMask = nsIpLineArray[1]
		snips = append(snips, snip)
	}
	return snips, nil
}

func ConvertMask(mask string) string {
	maskMap := SubnetMaskMap()
	decimalMask := maskMap[mask]
	return "/" + decimalMask
}

func GetNetworks(snips []Snip) ([]*net.IPNet, error) {
	var networks []*net.IPNet
	for _, snip := range snips {
		_, network, err := net.ParseCIDR(snip.ipAddress + ConvertMask(snip.subnetMask))
		if err != nil {
			return []*net.IPNet{}, err
		}
		networks = append(networks, network)
	}
	return networks, nil
}

func SubnetMaskMap() map[string]string {
	subnetMap := make(map[string]string)
	subnetMap["255.0.0.0"] = "8"
	subnetMap["255.128.0.0"] = "9"
	subnetMap["255.192.0.0"] = "10"
	subnetMap["255.224.0.0"] = "11"
	subnetMap["255.240.0.0"] = "12"
	subnetMap["255.248.0.0"] = "13"
	subnetMap["255.252.0.0"] = "14"
	subnetMap["255.254.0.0"] = "15"
	subnetMap["255.255.0.0"] = "16"
	subnetMap["255.255.128.0"] = "17"
	subnetMap["255.255.192.0"] = "18"
	subnetMap["255.255.224.0"] = "19"
	subnetMap["255.255.240.0"] = "20"
	subnetMap["255.255.248.0"] = "21"
	subnetMap["255.255.252.0"] = "22"
	subnetMap["255.255.254.0"] = "23"
	subnetMap["255.255.255.0"] = "24"
	subnetMap["255.255.255.128"] = "25"
	subnetMap["255.255.255.192"] = "26"
	subnetMap["255.255.255.224"] = "27"
	subnetMap["255.255.255.240"] = "28"
	subnetMap["255.255.255.248"] = "29"
	subnetMap["255.255.255.252"] = "30"
	subnetMap["255.255.255.254"] = "31"
	subnetMap["255.255.255.255"] = "32"
	return subnetMap
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println(os.Stderr, "Usage: %s filename\n", os.Args[0])
	}
	filename := os.Args[1]
	snips, err := GetSnips(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	networks, err := GetNetworks(snips)
	if err != nil {
		fmt.Println(err)
		return
	}
	servers, err := GetServers(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	serverMap := make(map[string]string)
	for _, network := range networks {
		for _, server := range servers {
			serverIP := net.ParseIP(server.ipAddress)
			networkCheck := network.Contains(serverIP)
			if networkCheck == true {
				serverMap[server.ipAddress] = serverMap[server.ipAddress]
				serverMap[server.ipAddress] = server.ipAddress
			}
		}
	}
	for _, server := range servers {
		if serverMap[server.ipAddress] != server.ipAddress {
			fmt.Println(server.ipAddress)
		}
	}
}
