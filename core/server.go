package core

import (
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"sephiroth-go/init"
	"sephiroth-go/service/sys"
	"time"
)

type server interface {
	ListenAndServe() error
}

func initServer(address string, router *gin.Engine) server {
	s := endless.NewServer(address, router)
	s.ReadHeaderTimeout = 20 * time.Second
	s.WriteTimeout = 20 * time.Second
	s.MaxHeaderBytes = 1 << 20
	return s
}

func RunWindowsServer() {
	if Config.System.UseMultipoint || Config.System.UseRedis {
		// 初始化redis服务
		init.Redis()
	}

	// 从db加载jwt数据
	if Db != nil {
		sys.LoadAll()
	}

	Router := init.Routers()
	Router.Static("/form-generator", "./resource/page")

	address := fmt.Sprintf(":%d", Config.System.Addr)
	s := initServer(address, Router)

	Log.Info("server run success on ", zap.String("address", address))

	fmt.Printf(`%s`, address)
	Log.Error(s.ListenAndServe().Error())
}
