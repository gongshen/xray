package business

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gongshen/xray/stat/conn"
	"github.com/gongshen/xray/stat/dao"
	"sort"
)

func InitStat() {
	if err := conn.InitDB(XrayDBPath); err != nil {
		panic(err)
	}
}

func GetStat(c *gin.Context) {
	traffics, err := dao.GetEnableTraffics()
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	sort.Slice(traffics, func(i, j int) bool {
		return traffics[i].Base+traffics[i].Used > traffics[j].Base+traffics[j].Used
	})
	var output bytes.Buffer
	output.WriteString("<html><body>")
	output.WriteString("---------------------------------------<br>")
	for _, traffic := range traffics {
		output.WriteString(fmt.Sprintf(`<font color="red">用户：</font>%s<br><font color="green">已用流量：</font>%s<br>`, traffic.Tag, format(traffic.Used+traffic.Base)))
		output.WriteString("---------------------------------------<br>")
	}
	output.WriteString("</body></html>")
	c.Data(200, "text/html; charset=utf-8", output.Bytes())
}

func format(used int64) string {
	if used < GB {
		return fmt.Sprintf("%.2fMB", float64(used)/MB)
	} else {
		return fmt.Sprintf("%.2fGB", float64(used)/GB)
	}
}
