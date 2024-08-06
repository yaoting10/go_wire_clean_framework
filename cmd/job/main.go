package main

import (
	"fmt"
	"github.com/gophero/goal/ginx"
	"github.com/gophero/goal/logx"
	"github.com/gophero/goal/mailx"
	"goboot/cmd/job/conf"
	"goboot/internal/config"
)

func main() {
	app, cleanup, err := newApp()
	if err != nil {
		panic(err)
	}
	c := app.Conf
	fmt.Printf("server start: %s:%d\n", "http://"+c.Http().Ip, c.Http().Port)

	ginx.Run(app.Engine, fmt.Sprintf(":%d", c.Http().Port))
	defer cleanup()
}

// 配置provider
func ConfProvider() *conf.JobConf {
	return conf.NewConfig(func(conf config.Conf) {
	}, "cmd/job/conf/dev.yml")
}
func MailConfProvider(c *conf.JobConf) mailx.Conf {
	return mailx.Conf{Server: c.MailConf.Server, SupportDomain: c.MailConf.SupportDomain}
}

// LoggerProvider provide Logger
func LoggerProvider(c *conf.JobConf) *logx.Logger {
	zap := c.Zap()
	return logx.NewLog(&zap)
}
