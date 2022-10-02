package business

import (
	"github.com/gin-gonic/gin"
	"github.com/gongshen/xray/dao"
	"github.com/gongshen/xray/models"
	"github.com/xtls/xray-core/common/net"
	"net/http"
)

type GinServer struct {
	l net.Listener
}

func NewGinServer() *GinServer {
	return &GinServer{}
}

func (g *GinServer) Start() error {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		return err
	}
	g.l = listener

	engine := gin.New()
	srv := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: engine,
	}
	engine.GET("/stat/get", GetTraffic)
	engine.GET("/stat/new", NewTraffic)

	go srv.Serve(listener)
	return nil
}

func (g *GinServer) Close() {
	g.l.Close()
}

func GetTraffic(c *gin.Context) {
	ans, err := dao.GetTraffics()
	if err != nil {
		c.JSON(500, err)
	} else {
		c.JSON(200, ans)
	}
}

func NewTraffic(c *gin.Context) {
	tag := c.Query("tag")
	if err := dao.SaveTraffic(&models.Traffic{
		Tag: tag,
	}); err != nil {
		c.JSON(500, err)
	} else {
		c.JSON(200, "OK")
	}
}
