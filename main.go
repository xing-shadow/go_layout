package main

import (
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	"layout/config"
	"layout/global"
	"layout/internal/serve"
	"layout/internal/service"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

var (
	ConfigPath string
)

func main() {
	a := kingpin.New(filepath.Base(os.Args[0]), " config")
	a.HelpFlag.Short('h')
	a.Flag("conf", "config file path").Short('c').StringVar(&ConfigPath)
	if _, err := a.Parse(os.Args[1:]); err != nil {
		panic("解析命令行参数失败: " + err.Error())
	}
	if err := config.Init(ConfigPath); err != nil {
		panic(err)
	}
	if err := global.Init(global.Option{LoggerCfg: &config.Cfg.Logger}); err != nil {
		panic(err)
	}
	if err := service.Init(); err != nil {
		panic(err)
	}
	var exit = make(chan error, 1)
	serve.Start(exit)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
		//写入日志
		global.GetLogger().Sync()
		//关闭服务
		serve.Stop()
		fmt.Println("程序退出")
	case err := <-exit:
		//服务启动错误
		panic("服务启动失败:" + err.Error())
	}
}
