package dialer

import (
	"net"
	"strconv"
	"syscall"
	"time"

	"github.com/slackhq/nebula/component/iface"
	"github.com/slackhq/nebula/util"
	"golang.org/x/sys/unix"
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

		ipStr, _, err := net.SplitHostPort(address)
		if err == nil {
			ip := net.ParseIP(ipStr)
			if ip != nil && !ip.IsGlobalUnicast() {
				return
			}
		}

		var innerErr error
		err = c.Control(func(fd uintptr) {
			switch network {
			case "tcp4", "udp4":
				innerErr = unix.SetsockoptInt(int(fd), unix.IPPROTO_IP, unix.IP_BOUND_IF, ifaceIdx)
			case "tcp6", "udp6":
				innerErr = unix.SetsockoptInt(int(fd), unix.IPPROTO_IPV6, unix.IPV6_BOUND_IF, ifaceIdx)
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
	// 开启监听
	go netInterfaceListener()
	return address, nil
}

func netInterfaceListener() {
	for {
		time.Sleep(2 * time.Second)
		_, indexStr, err := util.InitInterface()
		if err == nil {
			index, errInt := strconv.Atoi(indexStr)
			if errInt != nil {
				ifaceIdx = index
			}
		}
	}
}
