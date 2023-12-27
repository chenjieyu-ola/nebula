package util

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var InterfaceName string
var InterfaceIndex = 0

func InitInterface() (string, int, error) {
	if runtime.GOOS == "darwin" {
		netstatCmd := exec.Command("bash", "-c", "netstat -nr | grep default")
		netstatOutput, err := netstatCmd.Output()
		defaultInfName := ""
		if err != nil {
			fmt.Println(err)
			return "", 0, err
		} else {
			defaultInf := strings.Fields(strings.Split(string(netstatOutput), "\n")[0])
			defaultInfName = defaultInf[len(defaultInf)-1]
			fmt.Println(defaultInfName)
		}
		if defaultInfName == "" {
			return "", 0, fmt.Errorf("get default interfaces err")
		}
		inter, _ := net.Interfaces()
		for _, i := range inter {
			if i.Name == defaultInfName {
				InterfaceIndex = i.Index
				return defaultInfName, i.Index, nil
			}
		}

		return "", 0, fmt.Errorf("get default interfaces index err")

	}
	return "", 0, fmt.Errorf("get Interface error")
}

func NetInterfaceListener() {
	for {
		time.Sleep(2 * time.Second)
		_, index, err := InitInterface()
		if err == nil {
			InterfaceIndex = index
		}
	}
}
