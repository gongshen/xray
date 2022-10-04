package business

import (
	"encoding/json"
	"github.com/gongshen/xray/stat/models"
	"io/ioutil"
)

// GetConfigFromFile 读取配置文件
func GetConfigFromFile() (*models.XrayConfig, error) {
	// 读取文件
	confData, err := ioutil.ReadFile(XrayConfigFile)
	if err != nil {
		return nil, err
	}
	cnf := new(models.XrayConfig)
	if err = json.Unmarshal(confData, cnf); err != nil {
		return nil, err
	}
	return cnf, nil
}

// SaveConfigToFile 将配置写入文件（覆盖）
func SaveConfigToFile(cnf *models.XrayConfig) error {
	data, err := json.Marshal(cnf)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(XrayConfigFile, data, 0644)
}
