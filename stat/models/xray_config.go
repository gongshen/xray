package models

type XrayConfig struct {
	Policy          *XrayConfigPolicy     `json:"policy,omitempty"`
	Stats           XrayConfigStats       `json:"stats,omitempty"`
	LogConfig       *XrayConfigLog        `json:"log,omitempty"`
	RouterConfig    *XrayConfigRouter     `json:"routing,omitempty"`
	InboundConfigs  []*XrayConfigInbound  `json:"inbounds,omitempty"`
	OutboundConfigs []*XrayConfigOutbound `json:"outbounds,omitempty"`
	API             *XrayConfigApi        `json:"api,omitempty"`
}

type XrayConfigInbound struct {
	Listen         string                    `json:"listen,omitempty"`
	Port           int64                     `json:"port,omitempty"`
	Protocol       string                    `json:"protocol,omitempty"`
	Tag            string                    `json:"tag,omitempty"`
	Settings       *XrayConfigSettings       `json:"settings,omitempty"`
	StreamSettings *XrayConfigStreamSettings `json:"streamSettings,omitempty"`
}

type XrayConfigOutbound struct {
	Protocol       string                    `json:"protocol,omitempty"`
	Settings       *XrayConfigSettings       `json:"settings,omitempty"`
	StreamSettings *XrayConfigStreamSettings `json:"streamSettings,omitempty"`
}

type XrayConfigSettings struct {
	Address    string                        `json:"address,omitempty"`
	Decryption string                        `json:"decryption,omitempty"`
	Clients    []*XrayConfigSettingsClient   `json:"clients,omitempty"`
	Fallbacks  []*XrayConfigSettingsFallback `json:"fallbacks,omitempty"`
}

type XrayConfigStreamSettings struct {
	Network     string                 `json:"network,omitempty"`
	Security    string                 `json:"security,omitempty"`
	TlsSettings *XrayConfigTlsSettings `json:"tlsSettings,omitempty"`
	TcpSettings *XrayConfigTcpSettings `json:"tcpSettings,omitempty"`
}

type XrayConfigTcpSettings struct {
	Header *XrayConfigTcpSettingsHeader `json:"header,omitempty"`
}

type XrayConfigTcpSettingsHeader struct {
	Typ      string                               `json:"type,omitempty"`
	Response *XrayConfigTcpSettingsHeaderResponse `json:"response,omitempty"`
}

type XrayConfigTcpSettingsHeaderResponse struct {
	Version string                                      `json:"version,omitempty"`
	Status  string                                      `json:"status,omitempty"`
	Reason  string                                      `json:"reason,omitempty"`
	Headers *XrayConfigTcpSettingsHeaderResponseHeaders `json:"headers,omitempty"`
}

type XrayConfigTcpSettingsHeaderResponseHeaders struct {
	ContentType      []string `json:"Content-Type,omitempty"`
	TransferEncoding []string `json:"Transfer-Encoding,omitempty"`
	Connection       []string `json:"Connection,omitempty"`
	Pragma           string   `json:"pragma,omitempty"`
}

type XrayConfigTlsSettings struct {
	ServerName   string                              `json:"serverName,omitempty"`
	Alpn         []string                            `json:"alpn,omitempty"`
	Certificates []*XrayConfigTlsSettingsCertificate `json:"certificates,omitempty"`
}

type XrayConfigTlsSettingsCertificate struct {
	CertificateFile string `json:"certificateFile,omitempty"`
	KeyFile         string `json:"keyFile,omitempty"`
}

type XrayConfigSettingsFallback struct {
	Dest int64  `json:"dest,omitempty"`
	Alpn string `json:"alpn,omitempty"`
	Xver int64  `json:"xver,omitempty"`
}

type XrayConfigSettingsClient struct {
	Id      string `json:"id,omitempty"`
	Level   int64  `json:"level"`
	Email   string `json:"email,omitempty"`
	AlterId int64  `json:"alterId,omitempty"`
}

type XrayConfigApi struct {
	Tag      string   `json:"tag,omitempty"`
	Services []string `json:"services,omitempty"`
}

type XrayConfigRouter struct {
	DomainStrategy string                  `json:"domainStrategy,omitempty"`
	Rules          []*XrayConfigRouterRule `json:"rules,omitempty"`
}

type XrayConfigRouterRule struct {
	Typ         string   `json:"type,omitempty"`
	InboundTag  []string `json:"inboundTag,omitempty"`
	OutboundTag string   `json:"outboundTag,omitempty"`
	Ip          []string `json:"ip,omitempty"`
	Domain      []string `json:"domain,omitempty"`
	Protocol    []string `json:"protocol,omitempty"`
}

type XrayConfigLog struct {
	Access   string `json:"access,omitempty"`
	Error    string `json:"error,omitempty"`
	LogLevel string `json:"loglevel,omitempty"`
}

type XrayConfigStats struct{}

type XrayConfigPolicy struct {
	Levels map[string]*XrayConfigPolicyLevel `json:"levels,omitempty"`
	System *XrayConfigPolicySystem           `json:"system,omitempty"`
}

type XrayConfigPolicyLevel struct {
	StatsUserUplink   bool `json:"statsUserUplink,omitempty"`
	StatsUserDownlink bool `json:"statsUserDownlink,omitempty"`
}

type XrayConfigPolicySystem struct {
	StatsInboundUplink    bool `json:"statsInboundUplink,omitempty"`
	StatsInboundDownlink  bool `json:"statsInboundDownlink,omitempty"`
	StatsOutboundUplink   bool `json:"statsOutboundUplink,omitempty"`
	StatsOutboundDownlink bool `json:"statsOutboundDownlink,omitempty"`
}
