package core

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

const (
	Env               = "default"
	ConfigDefaultFile = "config.yaml"
	ConfigTestFile    = "config.test.yaml"
	ConfigDebugFile   = "config.debug.yaml"
	ConfigReleaseFile = "config.release.yaml"
)

// Viper //
// 优先级: 命令行 > 环境变量 > 默认值
func Viper(path ...string) *viper.Viper {
	var config string

	if len(path) == 0 {
		flag.StringVar(&config, "c", "", "choose config file.")
		flag.Parse()
		if config == "" { // 判断命令行参数是否为空
			if configEnv := os.Getenv(Env); configEnv == "" { // 判断 constant.ConfigEnv 常量存储的环境变量是否为空
				switch gin.Mode() {
				case gin.DebugMode:
					config = ConfigDebugFile
					fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.Mode(), ConfigDebugFile)
				case gin.ReleaseMode:
					config = ConfigReleaseFile
					fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.Mode(), ConfigReleaseFile)
				case gin.TestMode:
					config = ConfigTestFile
					fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.Mode(), ConfigTestFile)
				}
			} else { // constant.ConfigEnv 常量存储的环境变量不为空 将值赋值于config
				config = configEnv
				fmt.Printf("您正在使用%s环境变量,config的路径为%s\n", Env, config)
			}
		} else { // 命令行参数不为空 将值赋值于config
			fmt.Printf("您正在使用命令行的-c参数传递的值,config的路径为%s\n", config)
		}
	} else { // 函数传递的可变参数的第一个值赋值于config
		config = path[0]
		fmt.Printf("您正在使用func Viper()传递的值,config的路径为%s\n", config)
	}

	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err = v.Unmarshal(&Config); err != nil {
			fmt.Println(err)
		}
	})
	if err = v.Unmarshal(&Config); err != nil {
		panic(err)
	}

	// root 适配性 根据root位置去找到对应迁移位置,保证root路径有效
	Config.AutoCode.Root, _ = filepath.Abs("..")
	return v
}
