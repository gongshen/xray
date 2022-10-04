package business

import "regexp"

const (
	TrafficTarget = "127.0.0.1:11111"

	XrayDBPath     = "/usr/local/etc/xray.db"
	XrayConfigFile = "config.json"
	GB             = 1024 * 1024 * 1024
	MB             = 1024 * 1024
)

var (
	trafficRegex = regexp.MustCompile("(inbound|outbound|user)>>>([^>]+)>>>traffic>>>(downlink|uplink)")
)
