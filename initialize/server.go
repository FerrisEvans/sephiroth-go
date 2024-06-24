package initialize

import (
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"sephiroth-go/core"
	"sephiroth-go/core/log"
	system "sephiroth-go/service/sys"
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
	if core.Config.System.UseMultipoint || core.Config.System.UseRedis {
		// 初始化redis服务
		Redis()
	}

	// 从db加载jwt数据
	if core.Db != nil {
		system.LoadAll()
	}

	Router := Routers()
	Router.Static("/form-generator", "./resource/page")

	address := fmt.Sprintf(":%d", core.Config.System.Addr)
	s := initServer(address, Router)

	log.Log.Info("server run success on ", zap.String("address", address))

	fmt.Printf(`%s`, address)
	log.Log.Error(s.ListenAndServe().Error())
}
