package utils

import (
	"github.com/sirupsen/logrus"
	"net"
	"os"
)

var (
	Ip       string
	RemoteIp string
)

func SetRemoteIp() {
	RemoteIp = os.Getenv("REMOTE_IP")
}

func SetIp() {
	Ip = GetIp()
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
