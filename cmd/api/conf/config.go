package conf

import (
	"flag"
	"fmt"
	"github.com/gophero/goal/aws/s3"
	"github.com/gophero/goal/logx"
	"goboot/internal/config"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type ApiConf struct {
	HttpConf    config.Http    `mapstructure:"http" yaml:"http"`
	AssetConf   config.Asset   `mapstructure:"asset" yaml:"asset"`
	I18nConf    config.I18n    `mapstructure:"i18n" yaml:"i18n"`
	DBsConf     config.DBs     `mapstructure:",squash"`
	ZapConf     logx.Zap       `mapstructure:"zap" yaml:"zap"`
	PortalConf  config.Portal  `mapstructure:"portal" yaml:"portal"`
	AuthConf    config.Auth    `mapstructure:"auth" yaml:"auth"`
	RedisConf   config.Redis   `mapstructure:"redis" yaml:"redis"`
	MongoConf   config.Mongo   `mapstructure:"mongo" yaml:"mongo"`
	AwsS3Conf   config.AwsS3   `mapstructure:"aws-s3" yaml:"aws-s3"`
	GoogleConf  config.Google  `mapstructure:"google" yaml:"google"`
	MailConf    config.Mail    `mapstructure:"mail" yaml:"mail"`
	TwitterConf config.Twitter `mapstructure:"twitter" yaml:"twitter"`
	*viper.Viper
}

func (a *ApiConf) Http() config.Http {
	return a.HttpConf
}

func (a *ApiConf) Asset() config.Asset {
	return a.AssetConf
}

func (a *ApiConf) I18n() config.I18n {
	return a.I18nConf
}

func (a *ApiConf) Zap() logx.Zap {
	return a.ZapConf
}

func (a *ApiConf) Server() config.Server {
	return config.Server{AuthConf: a.AuthConf, Portal: a.PortalConf}
}

func (a *ApiConf) Storage() config.Storage {
	return config.Storage{
		DBsConf:   a.DBsConf,
		RedisConf: a.RedisConf,
		MongoConf: a.MongoConf,
	}
}

func (a *ApiConf) Job() config.Job {
	panic("not implemented")
}

func (a *ApiConf) AwsS3() config.AwsS3 {
	return a.AwsS3Conf
}

func (a *ApiConf) Google() config.Google {
	return a.GoogleConf
}

func (a *ApiConf) Mail() config.Mail {
	return a.MailConf
}

func (a *ApiConf) Twitter() config.Twitter {
	return a.TwitterConf
}

func (a *ApiConf) GetViper() *viper.Viper {
	return a.Viper
}

func (a *ApiConf) SetViper(v *viper.Viper) {
	a.Viper = v
}

func NewConfig(f func(conf config.Conf), def ...string) *ApiConf {
	envConf := os.Getenv("APP_CONF")
	var ip string
	var port int
	var logPath string
	flag.StringVar(&ip, "ip", "", "配置服务ip")
	flag.IntVar(&port, "p", 0, "配置服务运行端口")
	flag.StringVar(&logPath, "logPath", "", "配置日志存储路径")
	if envConf == "" {
		flag.StringVar(&envConf, "c", "", "config path, eg: -c config/conf.yml")
	}
	flag.Parse()
	if envConf == "" && len(def) > 0 {
		envConf = def[0]
	}
	if envConf == "" {
		panic("no config file to load, use \"-c\" option to set one.")
	}
	if strings.Contains(envConf, "dev.yml") {
		gin.SetMode(gin.DebugMode)
	}
	fmt.Println("load v file:", envConf)

	v := viper.New()
	v.SetConfigFile(envConf)
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var conf = &ApiConf{}

	// 解析配置文件到结构体
	if err = v.Unmarshal(&conf); err != nil {
		fmt.Println("parse config file error: ", err)
	}

	http := config.PreferHttp(config.Http{Ip: ip, Port: port}, conf)
	v.Set("http.ip", http.Ip)
	v.Set("http.port", http.Port)
	if logPath == "" {
		logPath = conf.Zap().Director
	}
	v.Set("zap.director", logPath)
	http.Env = conf.HttpConf.Env
	conf.HttpConf = http

	conf.SetViper(v)
	AfterConfigLoaded(conf, f)

	return conf
}

func NewS3Conf(conf config.Conf) *s3.Conf {
	return &s3.Conf{Region: conf.AwsS3().Region, AccessKey: conf.AwsS3().AccessKeyId, AccessSecret: conf.AwsS3().AccessKeySecret}
}

func AfterConfigLoaded(conf config.Conf, f func(conf config.Conf)) {
	f(conf)
}
