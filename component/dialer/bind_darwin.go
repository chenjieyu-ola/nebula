package dialer

import (
	"fmt"
	"github.com/slackhq/nebula/component/iface"
	"github.com/slackhq/nebula/util"
	"golang.org/x/sys/unix"
	"net"
	"syscall"
)

type controlFn = func(network, address string, c syscall.RawConn) error

var ifaceIdx = 0

func bindControl(chain controlFn) controlFn {
	return func(network, address string, c syscall.RawConn) (err error) {
		defer func() {
			if err == nil && chain != nil {
				err = chain(network, address, c)
			}
		}()

		//ipStr, _, err := net.SplitHostPort(address)
		//if err == nil {
		//	ip := net.ParseIP(ipStr)
		//	if ip != nil && !ip.IsGlobalUnicast() {
		//		fmt.Println("=====================================")
		//		return
		//	}
		//}

		var innerErr error
		err = c.Control(func(fd uintptr) {
			fmt.Println("=====================================", fd)
			switch network {
			case "tcp4", "udp4":
				fmt.Println("=====================================", util.InterfaceIndex)
				innerErr = unix.SetsockoptInt(int(fd), unix.IPPROTO_IP, unix.IP_BOUND_IF, util.InterfaceIndex)
			case "tcp6", "udp6":
				innerErr = unix.SetsockoptInt(int(fd), unix.IPPROTO_IPV6, unix.IPV6_BOUND_IF, util.InterfaceIndex)
			}
		})

		if innerErr != nil {
			err = innerErr
		}

		return
	}
}

func bindIfaceToDialer(ifaceName string, dialer *net.Dialer, _ string, _ net.IP) error {
	ifaceObj, err := iface.ResolveInterface(ifaceName)
	if err != nil {
		return err
	}
	ifaceIdx = ifaceObj.Index
	dialer.Control = bindControl(dialer.Control)
	return nil
}

func BindIfaceToListenConfig(ifaceName string, lc *net.ListenConfig, _, address string) (string, error) {
	ifaceObj, err := iface.ResolveInterface(ifaceName)
	if err != nil {
		return "", err
	}
	ifaceIdx = ifaceObj.Index
	lc.Control = bindControl(lc.Control)
	return address, nil
}
