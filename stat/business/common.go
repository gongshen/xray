package business

const (
	VLessPattern   = `vless://{uid}@{ip}:{port}?security=tls#dino1-{ip}`
	XrayConfigFile = "/usr/local/etc/xray/config2.json"
)

var testContent = `
{
  "stat": [
    {
      "name": "inbound>>>vmess-quic>>>traffic>>>downlink",
      "value": 1176
    },
    {
      "name": "user>>>1>>>traffic>>>downlink",
      "value": 2040
    },
    {
      "name": "inbound>>>api>>>traffic>>>uplink",
      "value": 14247
    },
    {
      "name": "user>>>2>>>traffic>>>uplink",
      "value": 2520
    },
 	{
      "name": "user>>>3>>>traffic>>>uplink",
      "value": 252000
    },
    {
      "name": "inbound>>>api>>>traffic>>>downlink",
      "value": 87618
    },
    {
      "name": "outbound>>>direct>>>traffic>>>downlink",
      "value": 0
    },
    {
      "name": "inbound>>>vmess-quic>>>traffic>>>uplink",
      "value": 1691
    },
    {
      "name": "outbound>>>direct>>>traffic>>>uplink",
      "value": 0
    }
  ]
}
`
