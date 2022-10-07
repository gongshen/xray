package business

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	ginServer *GinServer
	jobServer *JobServer
}

func NewServer() *Server {
	return &Server{
		ginServer: NewGinServer(),
		jobServer: NewJobServer(),
	}
}

func (s *Server) Start() {
	if err := s.jobServer.Start(); err != nil {
		logrus.Println(err)
		return
	}
	if err := s.ginServer.Start(); err != nil {
		logrus.Println(err)
		return
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
	s.Close()
	time.Sleep(1 * time.Second)
}

func (s *Server) Close() {
	s.jobServer.Close()
	s.ginServer.Close()
}
