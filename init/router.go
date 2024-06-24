package init

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"sephiroth-go/core"
	"sephiroth-go/middleware"
)

type justFilesFilesystem struct {
	fs http.FileSystem
}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}

	stat, err := f.Stat()
	if stat.IsDir() {
		return nil, os.ErrPermission
	}

	return f, nil
}

func Routers() *gin.Engine {
	Router := gin.New()
	Router.Use(gin.Recovery())
	if gin.Mode() == gin.DebugMode {
		Router.Use(gin.Logger())
	}

	Router.StaticFS(core.Config.Local.StorePath, justFilesFilesystem{http.Dir(core.Config.Local.StorePath)}) // Router.Use(middleware.LoadTls())  // 如果需要使用https 请打开此中间件 然后前往 core/server.go 将启动模式 更变为 Router.RunTLS("端口","你的cre/pem文件","你的key文件")
	// 跨域，如需跨域可以打开下面的注释
	// Router.Use(middleware.Cors()) // 直接放行全部跨域请求
	// Router.Use(middleware.CorsByRules()) // 按照配置的规则放行跨域请求
	//core.Log.Info("use middleware cors")

	// 方便统一添加路由组前缀 多服务器上线使用

	PublicGroup := Router.Group(core.Config.System.RouterPrefix)

	PrivateGroup := Router.Group(core.Config.System.RouterPrefix)

	PrivateGroup.Use(middleware.JwtAuth()).Use(middleware.CasbinHandler())

	{
		// 健康监测
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})
	}

	//authRouter := router.RouterGroupApp.Auth
	//exampleRouter := router.RouterGroupApp.Storage
	//sysRouter := router.RouterGroupApp.Sys
	{
		//systemRouter.InitBaseRouter(PublicGroup) // 注册基础功能路由 不做鉴权
		//systemRouter.InitInitRouter(PublicGroup) // 自动初始化相关
	}

	{
		//systemRouter.InitApiRouter(PrivateGroup, PublicGroup)       // 注册功能api路由
		//systemRouter.InitJwtRouter(PrivateGroup)                    // jwt相关路由
		//systemRouter.InitUserRouter(PrivateGroup)                   // 注册用户路由
		//systemRouter.InitMenuRouter(PrivateGroup)                   // 注册menu路由
		//systemRouter.InitSystemRouter(PrivateGroup)                 // system相关路由
		//systemRouter.InitCasbinRouter(PrivateGroup)                 // 权限相关路由
		//systemRouter.InitAutoCodeRouter(PrivateGroup)               // 创建自动化代码
		//systemRouter.InitAuthorityRouter(PrivateGroup)              // 注册角色路由
		//systemRouter.InitSysDictionaryRouter(PrivateGroup)          // 字典管理
		//systemRouter.InitAutoCodeHistoryRouter(PrivateGroup)        // 自动化代码历史
		//systemRouter.InitSysOperationRecordRouter(PrivateGroup)     // 操作记录
		//systemRouter.InitSysDictionaryDetailRouter(PrivateGroup)    // 字典详情管理
		//systemRouter.InitAuthorityBtnRouterRouter(PrivateGroup)     // 字典详情管理
		//systemRouter.InitSysExportTemplateRouter(PrivateGroup)      // 导出模板
		//exampleRouter.InitCustomerRouter(PrivateGroup)              // 客户路由
		//exampleRouter.InitFileUploadAndDownloadRouter(PrivateGroup) // 文件上传下载功能路由

	}

	//插件路由安装
	//InstallPlugin(PrivateGroup, PublicGroup)

	// 注册业务路由
	//initBizRouter(PrivateGroup, PublicGroup)

	core.Log.Info("router register success")
	return Router
}
