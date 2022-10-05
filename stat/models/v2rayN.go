package models

type V2rayN struct {
	V    string `json:"v"`    //配置文件版本号,主要用来识别当前配置
	Ps   string `json:"ps"`   //备注或别名
	Add  string `json:"add"`  //地址IP或域名
	Port string `json:"port"` //端口号
	Id   string `json:"id"`   //UUID
	Aid  string `json:"aid"`  //alterId
	Scy  string `json:"scy"`  //加密方式(security),没有时值默认auto
	Net  string `json:"net"`  //传输协议(tcp\kcp\ws\h2\quic)
	Type string `json:"type"` //伪装类型(none\http\srtp\utp\wechat-video) *tcp or kcp or QUIC
	Host string `json:"host"` //伪装的域名
	Path string `json:"path"`
	Tls  string `json:"tls"` //底层传输安全(tls)
	Sni  string `json:"sni"` //serverName
}
