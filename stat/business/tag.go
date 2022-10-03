package business

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/gongshen/xray/stat/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func NewTag(c *gin.Context) {
	tag := c.Query("tag")
	if tag == "" {
		c.JSON(500, errors.New("参数错误"))
		return
	}
	// 读取文件
	confData, err := ioutil.ReadFile(XrayConfigFile)
	if err != nil {
		c.JSON(500, err)
		return
	}
	cnf := new(models.XrayConfig)
	if err = json.Unmarshal(confData, cnf); err != nil {
		fmt.Println(err)
		c.JSON(500, err)
		return
	}

	// 生成uuid
	id := uuid.New().String()
	newCli := &models.XrayConfigSettingsClient{
		Email: tag,
		Level: 0,
		Id:    id,
	}
	// 找到inbound中，listen为0.0.0.0的配置
	for _, inbound := range cnf.InboundConfigs {
		if inbound.Listen == "0.0.0.0" {
			// 遍历是否已存在tag
			for _, cli := range inbound.Settings.Clients {
				if cli.Email == tag {
					c.JSON(500, errors.New("tag已存在"))
					return
				}
			}
			inbound.Settings.Clients = append(inbound.Settings.Clients, newCli)
			break
		}
	}
	// 写入文件
	data, err := json.Marshal(cnf)
	if err != nil {
		c.JSON(500, err)
		return
	}
	if err = ioutil.WriteFile(XrayConfigFile, data, 0644); err != nil {
		c.JSON(500, err)
		return
	}
	if err = XrayRestart(); err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, "OK")
}

func DelTag(c *gin.Context) {
	tag := c.Query("tag")
	if tag == "" {
		c.JSON(500, errors.New("参数错误"))
		return
	}
	// 打开文件
	f, err := os.OpenFile(XrayConfigFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		c.JSON(500, err)
		return
	}
	defer f.Close()
	// 读取文件
	cnf := new(models.XrayConfig)
	confData, _ := ioutil.ReadAll(f)
	if err = json.Unmarshal(confData, cnf); err != nil {
		c.JSON(500, err)
		return
	}
	// 找到inbound中，listen为0.0.0.0的配置
	for _, inbound := range cnf.InboundConfigs {
		if inbound.Listen == "0.0.0.0" {
			newClients := make([]*models.XrayConfigSettingsClient, 0)
			// 遍历是否已存在tag
			for _, cli := range inbound.Settings.Clients {
				if cli.Email != tag {
					newClients = append(newClients, cli)
				} else {
					// 找到删除的tag
					logrus.Debugln("找到需要删除的tag")
				}
			}
			inbound.Settings.Clients = newClients
			break
		}
	}
	// 序列化
	data, err := json.Marshal(cnf)
	if err != nil {
		c.JSON(500, err)
		return
	}

	if err = ioutil.WriteFile(XrayConfigFile, data, 0644); err != nil {
		c.JSON(500, err)
		return
	}
	if err = XrayRestart(); err != nil {
		c.JSON(500, err)
		return
	}
	// 返回
	c.JSON(200, "OK")
}
