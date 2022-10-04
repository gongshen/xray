package business

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gongshen/xray/stat/models"
	"github.com/gongshen/xray/stat/utils"
)

const (
	VMessPattern = `vmess://tcp:%s-%d@%s:%d?type=http#dino2-%s`
	VLessPattern = `vless://%s@%s:%d?security=tls#dino1-%s"`

	VMessQrPattern = `https://api.qrserver.com/v1/create-qr-code/?size=400x400&data=vmess://tcp:%s-%d@%s:%d?type=http#dino2-%s`
	VLessQrPattern = `https://api.qrserver.com/v1/create-qr-code/?size=400x400&data=vless://%s@%s:%d?security=tls#dino1-%s`
)

func Share(c *gin.Context) {
	tag := c.Query("name")
	if tag == "" {
		c.JSON(500, "参数错误")
		return
	}
	cnf, err := GetConfigFromFile()
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	link, url, err := generateShare(tag, cnf)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	var output bytes.Buffer
	output.WriteString("<html><body><p>")
	output.WriteString(fmt.Sprintf(`<font color="red">用户：</font>%s<br><font color="green">分享链接：</font>%s<br><font color="green">二维码链接：</font>%s<br>`, tag, link, url))
	output.WriteString("</html></body></p>")
	c.Data(200, "text/html; charset=utf-8", output.Bytes())
}

func generateShare(tag string, cnf *models.XrayConfig) (string, string, error) {
	var (
		link string
		url  string
		err  error
		ip   = utils.GetIp()
	)
	for _, inbound := range cnf.InboundConfigs {
		if inbound.Listen == MainInboundListen {
			for _, client := range inbound.Settings.Clients {
				// 找到tag
				if client.Email == tag {
					if inbound.Protocol == "vless" {
						link = fmt.Sprintf(VLessPattern, client.Id, utils.Domain, inbound.Port, utils.Domain)
						url = fmt.Sprintf(VLessQrPattern, client.Id, utils.Domain, inbound.Port, utils.Domain)
					} else if inbound.Protocol == "vmess" {
						link = fmt.Sprintf(VMessPattern, client.Id, client.AlterId, ip, inbound.Port, ip)
						url = fmt.Sprintf(VMessQrPattern, client.Id, client.AlterId, ip, inbound.Port, ip)
					} else {
						err = errors.New("未知协议")
					}
					goto Find
				}
			}
			goto Find
		}
	}
Find:
	return link, url, err
}
