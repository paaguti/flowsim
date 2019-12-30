package common

import (
	"fmt"
	// 	"os"
	// 	"syscall"
	"errors"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"net"
)

func SetTos(conn net.Conn, dscp int) error {
	err := ipv4.NewConn(conn).SetTOS(dscp)
	if err != nil {
		// common.WarnErrorf(err, "while setting TOS")
		err = ipv6.NewConn(conn).SetTrafficClass(dscp)
	}
	return err
}

func SetTcpTos(conn *net.TCPConn, dscp int) error {
	err := ipv4.NewConn(conn).SetTOS(dscp)
	if err != nil {
		// common.WarnErrorf(err, "while setting TOS")
		err = ipv6.NewConn(conn).SetTrafficClass(dscp)
	}
	return err
}

func SetUdpTos(conn *net.UDPConn, dscp int) error {
	err := ipv4.NewConn(conn).SetTOS(dscp)
	if err != nil {
		// common.WarnErrorf(err, "while setting TOS")
		err = ipv6.NewConn(conn).SetTrafficClass(dscp)
	}

	return err
}

/*
 * Check whether the IP destination is IPv4 or IPv6
 * and set the UDP family to 'udp4' or 'udp6'
 */
func UdpFamily(ip string) (string, error) {
	ipAddr, err := net.ResolveIPAddr("ip", ip)
	if err == nil {
		if ipAddr.IP.To4() == nil {
			return "udp6", nil
		}
		return "udp4", nil
	}
	return "", err
}

func FirstIP(host string, ipv6 bool) (string, error) {
	//
	// Some systems name localhost differently for IPv6
	// Always use ::1 or 127.0.0.1 for localhost
	//
	if host == "localhost" {
		if ipv6 {
			return "::1", nil
		} else {
			return "127.0.0.1", nil
		}
	}
	family := "IPv4"
	if ipv6 {
		family = "IPv6"
	}
	ips, err := net.LookupIP(host)
	if err == nil {
		for _, ip := range ips {
			if ip.To4() == nil {
				if !ipv6 {
					continue
				}
			} else {
				if ipv6 {
					continue
				}
			}
			return ip.String(), nil
		}
		return "", errors.New(fmt.Sprintf("Couldn't find %s address for %s\n", family, host))
	}
	return "", err
}
