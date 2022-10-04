package utils

import (
	"github.com/sirupsen/logrus"
	"net"
)

var Domain string

func SetDomain(domain string) {
	Domain = domain
}

func GetIp() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		logrus.Errorln(err)
	}
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}
	}
	return ""
}
