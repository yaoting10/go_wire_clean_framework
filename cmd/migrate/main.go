package main

import (
	"github.com/gophero/goal/logx"
	"goboot/cmd/migrate/conf"
	"goboot/internal/config"
)

func main() {
	app, cleanup, err := newApp()
	if err != nil {
		panic(err)
	}
	app.Run()
	defer cleanup()
}

// 配置provider
func ConfProvider() *conf.MigrateConf {
	c := conf.NewConfig(func(conf config.Conf) {
	}, "cmd/migrate/conf/dev.yml")
	return &c
}

// LoggerProvider provide Logger
func LoggerProvider(c *conf.MigrateConf) *logx.Logger {
	zap := c.Zap()
	return logx.NewLog(&zap)
}
