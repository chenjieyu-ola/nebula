package util

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

var InterfaceName string
var ProxyInterface string

func InitInterface() (string, string, error) {

	var netIp = ""
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println(err)
	} else {
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		netIp = strings.Split(localAddr.String(), ":")[0]
		conn.Close()
	}

	//获取当前的所有设备
	inter, _ := net.Interfaces()
	var interList [][]string
	for _, i := range inter {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}
				// 解析出 Flags 包含有 up 属性的网卡
				// 对比网卡
				if strings.Contains(ip.String(), ".") && i.Flags&(1<<uint(0)) != 0 && i.HardwareAddr != nil && len(i.HardwareAddr) != 0 {
					fmt.Println("real ip ", ip.String(), "tun name ", i.Name, " mac ", i.HardwareAddr)
					//return ip.String(), terDevNames[ip.String()], i.HardwareAddr.String()
					interList = append(interList, []string{ip.String(), i.Name, strconv.Itoa(i.Index), i.HardwareAddr.String()})
				}
			}
		}
	}
	if len(interList) != 0 {
		for _, item := range interList {
			if item[0] == netIp {
				return item[1], item[2], nil
			}
		}
		// 取最后一个
		intef := interList[len(interList)-1]
		return intef[1], intef[2], nil

	}

	return "", "", fmt.Errorf("get Interface error")
}
