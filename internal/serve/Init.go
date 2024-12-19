package serve

import (
	"github.com/pkg/errors"
	"layout/config"
	"layout/internal/serve/httpServe"
)

func Start(errChan chan error) {
	go func() {
		err := httpServe.Init(&httpServe.Option{
			Port:      config.Cfg.HTTP.Port,
			ReadTime:  config.Cfg.HTTP.ReadTimeout,
			WriteTime: config.Cfg.HTTP.WriteTimeout,
		})
		if err != nil {
			err = errors.Wrap(err, "初始化http服务失败")
			errChan <- err
			return
		}
		err = httpServe.Start()
		if err != nil {
			err = errors.Wrap(err, "启动http服务失败")
			errChan <- err
			return
		}
	}()

	return
}

func Stop() {
	httpServe.Stop()
}
