package business

import (
	"github.com/gin-gonic/gin"
	"github.com/gongshen/xray/stat/dao"
	"github.com/gongshen/xray/stat/models"
	"github.com/google/uuid"
)

const (
	MainInboundListen = "0.0.0.0"
)

func NewUser(c *gin.Context) {
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
	// 生成uuid
	id := uuid.New().String()
	newCli := &models.XrayConfigSettingsClient{
		Email:   tag,
		Level:   0,
		Id:      id,
		AlterId: AlterID,
	}
	// 找到inbound中，listen为0.0.0.0的配置
	for _, inbound := range cnf.InboundConfigs {
		if inbound.Listen == MainInboundListen {
			// 遍历是否已存在tag
			for _, cli := range inbound.Settings.Clients {
				if cli.Email == tag {
					c.JSON(500, "用户已存在")
					return
				}
			}
			inbound.Settings.Clients = append(inbound.Settings.Clients, newCli)
			break
		}
	}
	if err = SaveConfigToFile(cnf); err != nil {
		c.JSON(500, err.Error())
		return
	}
	if err = Systemctl(SystemctlRestartOpt, ServiceNameXray); err != nil {
		c.JSON(500, err.Error())
		return
	}
	// 新建统计配置
	if err = dao.NewTraffic(tag); err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, "OK")
}

func DelUser(c *gin.Context) {
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
	find := false
	// 找到inbound中，listen为0.0.0.0的配置
	for _, inbound := range cnf.InboundConfigs {
		if inbound.Listen == MainInboundListen {
			newClients := make([]*models.XrayConfigSettingsClient, 0)
			// 遍历是否已存在tag
			for _, cli := range inbound.Settings.Clients {
				if cli.Email != tag {
					newClients = append(newClients, cli)
				} else {
					// 找到删除的tag
					find = true
				}
			}
			inbound.Settings.Clients = newClients
			break
		}
	}
	if !find {
		c.JSON(500, "未找到该用户")
		return
	}
	if err = SaveConfigToFile(cnf); err != nil {
		c.JSON(500, err.Error())
		return
	}
	if err = Systemctl(SystemctlRestartOpt, ServiceNameXray); err != nil {
		c.JSON(500, err.Error())
		return
	}
	// 软删除记录
	if err = dao.DelTraffic(tag); err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, "OK")
}
