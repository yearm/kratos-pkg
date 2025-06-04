package net

import (
	"github.com/yearm/kratos-pkg/errors"
	"net"
)

// GetIPArray retrieves and returns all the ip of current host.
func GetIPArray() ([]string, error) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		return nil, errors.Wrap(err, "net.InterfaceAddrs failed")
	}

	ips := make([]string, 0)
	for _, address := range interfaceAddr {
		ipNet, isValidIpNet := address.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips, nil
}

// GetFreePort retrieves and returns a port that is free.
func GetFreePort() (int, error) {
	var (
		network = `tcp`
		address = `:0`
	)
	resolvedAddr, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		return 0, errors.Wrap(err, "net.ResolveTCPAddr failed")
	}
	l, err := net.ListenTCP(network, resolvedAddr)
	if err != nil {
		return 0, errors.Wrap(err, "net.ListenTCP failed")
	}
	port := l.Addr().(*net.TCPAddr).Port
	if err = l.Close(); err != nil {
		return 0, errors.Wrap(err, "close listening failed")
	}
	return port, nil
}
