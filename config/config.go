package config

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"os"
)

var Cfg Config

type Config struct {
	Logger Logger `toml:"Logger"`
	HTTP   HTTP   `toml:"httpServe"`
}
type Logger struct {
	LogDir     string `toml:"LogDir"`
	FileName   string `toml:"FileName"`
	Level      string `toml:"Level"`
	IsProd     bool   `toml:"IsProd"`
	MaxSize    int    `toml:"MaxSize"`
	MaxBackups int    `toml:"MaxBackups"`
	MaxAge     int    `toml:"MaxAge"`
}
type HTTP struct {
	Port         int `toml:"Port"`
	ReadTimeout  int `toml:"ReadTimeout"`
	WriteTimeout int `toml:"WriteTimeout"`
}

func Init(path string) error {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "读取配置文件错误")
	}
	err = toml.Unmarshal(fileData, &Cfg)
	if err != nil {
		return errors.Wrap(err, "解析配置文件错误")
	}
	return nil
}
