package business

import (
	"github.com/gin-gonic/gin"
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
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}
	g.l = listener

	engine := gin.New()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}
	// 获取流量统计信息
	engine.GET("/user/stat", GetStat)
	engine.GET("/user/new", NewUser)
	engine.GET("/user/del", DelUser)
	engine.GET("/user/share", Share)
	engine.GET("/user/reset_stat", ResetStat)
	go srv.Serve(listener)
	return nil
}

func (g *GinServer) Close() {
	g.l.Close()
}

func ResetStat(c *gin.Context) {}
