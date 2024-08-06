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

type JobConf struct {
	HttpConf    config.Http    `mapstructure:"http" yaml:"http"`
	AssetConf   config.Asset   `mapstructure:"asset" yaml:"asset"`
	DBsConf     config.DBs     `mapstructure:",squash"`
	JobConf     config.Job     `mapstructure:"job" yaml:"job"`
	ZapConf     logx.Zap       `mapstructure:"zap" yaml:"zap"`
	RedisConf   config.Redis   `mapstructure:"redis" yaml:"redis"`
	MongoConf   config.Mongo   `mapstructure:"mongo" yaml:"mongo"`
	AwsS3Conf   config.AwsS3   `mapstructure:"aws-s3" yaml:"aws-s3"`
	MailConf    config.Mail    `mapstructure:"mail" yaml:"mail"`
	TwitterConf config.Twitter `mapstructure:"twitter" yaml:"twitter"`
	*viper.Viper
}

func (a JobConf) Http() config.Http {
	return a.HttpConf
}

func (a JobConf) Asset() config.Asset {
	return a.AssetConf
}

func (a JobConf) Zap() logx.Zap {
	return a.ZapConf
}

func (a JobConf) Server() config.Server {
	panic("not implemented")
}

func (a JobConf) Storage() config.Storage {
	return config.Storage{
		DBsConf:   a.DBsConf,
		RedisConf: a.RedisConf,
		MongoConf: a.MongoConf,
	}
}

func (a *JobConf) Job() config.Job {
	return a.JobConf
}

func (a *JobConf) I18n() config.I18n {
	panic("not implemented")
}

func (a *JobConf) AwsS3() config.AwsS3 {
	return a.AwsS3Conf
}

func (a *JobConf) Google() config.Google {
	panic("not implemented")
}

func (a *JobConf) Mail() config.Mail {
	return a.MailConf
}

func (a *JobConf) Twitter() config.Twitter {
	return a.TwitterConf
}

func (a *JobConf) GetViper() *viper.Viper {
	return a.Viper
}

func (a *JobConf) SetViper(v *viper.Viper) {
	a.Viper = v
}

func NewConfig(f func(conf config.Conf), def ...string) *JobConf {
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
	if strings.Contains(envConf, "dev.yaml") {
		gin.SetMode(gin.DebugMode)
	}
	fmt.Println("load v file:", envConf)

	v := viper.New()
	v.SetConfigFile(envConf)
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var conf = &JobConf{}

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
