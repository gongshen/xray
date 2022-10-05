package business

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gongshen/xray/stat/models"
	"github.com/gongshen/xray/stat/utils"
	"strconv"
)

const (
	VMessPattern = `vmess://tcp:%s-%d@%s:%d?type=http#dino2-%s`
	VLessPattern = `vless://%s@%s:%d?security=tls#dino1-%s`
	QrPattern    = `https://api.qrserver.com/v1/create-qr-code/?size=400x400&data=%s`
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
	shareInfo, err := shareFind(tag, cnf)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	var output bytes.Buffer
	output.WriteString("<html><body>")
	output.WriteString(fmt.Sprintf(`<font color="red">用户：</font>%s<br>`, tag))
	// shadowrocket,qv2ray
	output.WriteString("-----------------shadowrocket,qv2ray,v2rayXS---------------<br>")
	output.WriteString(fmt.Sprintf(`<font color="green">分享链接：</font>%s<br><font color="green">二维码：</font>`, shareInfo.Share1Link))
	if shareInfo.Share1Link != "" {
		output.WriteString(fmt.Sprintf(`<img src="%s" width="200" height="200"><br>`, fmt.Sprintf(QrPattern, shareInfo.Share1Link)))
	}
	output.WriteString("<br>")
	// v2rayN
	output.WriteString("-----------------v2rayN,v2rayXS---------------<br>")
	output.WriteString(fmt.Sprintf(`<font color="green">分享链接：</font>%s<br><font color="green">二维码：</font>`, shareInfo.Share2Link))
	if shareInfo.Share2Link != "" {
		output.WriteString(fmt.Sprintf(`<img src="%s" width="200" height="200"><br>`, fmt.Sprintf(QrPattern, shareInfo.Share2Link)))
	}
	output.WriteString("<br>")
	output.WriteString("</body></html>")
	c.Data(200, "text/html; charset=utf-8", output.Bytes())
}

func shareFind(tag string, cnf *models.XrayConfig) (*models.ShareInfo, error) {
	for _, inbound := range cnf.InboundConfigs {
		if inbound.Listen == MainInboundListen {
			for _, client := range inbound.Settings.Clients {
				// 找到tag
				if client.Email == tag {
					return generateShare(inbound.Protocol, inbound.Port, client)
				}
			}
		}
	}
	return nil, errors.New("未找到用户")
}

// generateShare 生成分享的链接
func generateShare(protocol string, port int64, client *models.XrayConfigSettingsClient) (*models.ShareInfo, error) {
	var (
		info = new(models.ShareInfo)
	)
	if protocol == "vless" {
		// share1
		info.Share1Link = fmt.Sprintf(VLessPattern, client.Id, utils.Domain, port, utils.Domain)
		// share2
		info.Share2Link = fmt.Sprintf(VLessPattern, client.Id, utils.Domain, port, utils.Domain)
	} else if protocol == "vmess" {
		// share1
		info.Share1Link = fmt.Sprintf(VMessPattern, client.Id, client.AlterId, utils.Ip, port, utils.Ip)
		// share2
		v2rayN := defaultV2rayN
		v2rayN.Add = utils.Ip
		v2rayN.Id = client.Id
		v2rayN.Port = strconv.Itoa(int(port))
		v2rayN.Ps = fmt.Sprintf("dino2-%s", utils.Ip)
		data, err := json.Marshal(v2rayN)
		if err != nil {
			return nil, err
		}
		info.Share2Link = fmt.Sprintf("vmess://%s", base64.StdEncoding.EncodeToString(data))
	} else {
		return nil, errors.New("未知协议")
	}
	return info, nil
}

var defaultV2rayN = models.V2rayN{
	V:    "2",
	Type: "http",
	Aid:  "64",
	Net:  "tcp",
	Scy:  "auto",
}
