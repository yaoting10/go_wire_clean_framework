package config

import (
	"github.com/gophero/goal/logx"
	"github.com/spf13/viper"
)

type Conf interface {
	Http() Http
	Asset() Asset
	I18n() I18n
	Zap() logx.Zap
	Server() Server
	Storage() Storage
	Job() Job
	AwsS3() AwsS3
	Google() Google
	Mail() Mail
	Twitter() Twitter

	GetViper() *viper.Viper
	SetViper(v *viper.Viper)
}

type Server struct {
	SignConf Sign   `mapstructure:"sign" yaml:"sign"`
	AuthConf Auth   `mapstructure:"auth" yaml:"auth"`
	Portal   Portal `mapstructure:"portal" yaml:"portal"`
}

type Storage struct {
	DBsConf   DBs   `mapstructure:"db" yaml:"db"`
	RedisConf Redis `mapstructure:"redis" yaml:"redis"`
	MongoConf Mongo `mapstructure:"mongo" yaml:"mongo"`
}

// PreferHttp 优先使用preferHttp提供的配置
func PreferHttp(preferHttp Http, conf Conf) Http {
	if preferHttp.Ip == "" {
		preferHttp.Ip = conf.Http().Ip
	}
	if preferHttp.Port == 0 {
		preferHttp.Port = conf.Http().Port
	}
	if preferHttp.Env == "" {
		preferHttp.Env = conf.Http().Env
	}
	return preferHttp
}
