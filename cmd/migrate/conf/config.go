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

type MigrateConf struct {
	DBsConf config.DBs `mapstructure:",squash"`
	ZapConf logx.Zap   `mapstructure:"zap" yaml:"zap"`
	*viper.Viper
}

func (a MigrateConf) Asset() config.Asset {
	panic("implement me")
}

func (a MigrateConf) Http() config.Http {
	panic("not implemented")
}

func (a MigrateConf) Zap() logx.Zap {
	return a.ZapConf
}

func (a MigrateConf) Server() config.Server {
	panic("not implemented")
}

func (a MigrateConf) Storage() config.Storage {
	return config.Storage{
		DBsConf:   a.DBsConf,
		RedisConf: config.Redis{},
		MongoConf: config.Mongo{},
	}
}

func (a MigrateConf) Job() config.Job {
	panic("not implemented")
}

func (a MigrateConf) I18n() config.I18n {
	panic("not implemented")
}

func (a MigrateConf) AwsS3() config.AwsS3 {
	panic("not implemented")
}

func (a MigrateConf) Google() config.Google {
	panic("not implemented")
}

func (a MigrateConf) Mail() config.Mail {
	panic("not implemented")
}

func (a MigrateConf) Twitter() config.Twitter {
	panic("not implemented")
}

func (a MigrateConf) GetViper() *viper.Viper {
	return a.Viper
}

func (a MigrateConf) SetViper(v *viper.Viper) {
	a.Viper = v
}

func NewConfig(f func(conf config.Conf), def ...string) MigrateConf {
	envConf := os.Getenv("APP_CONF")
	var ip string
	var port int
	var logPath string
	flag.StringVar(&ip, "ip", "127.0.0.1", "配置服务ip")
	flag.IntVar(&port, "p", 8080, "配置服务运行端口")
	flag.StringVar(&logPath, "logPath", "", "配置日志存储路径")
	if envConf == "" {
		flag.StringVar(&envConf, "c", "", "config path, eg: -c config/conf.yml")
		flag.Parse()
	}
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
	if ip != "" {
		v.Set("http.ip", ip)
	}
	if port > 0 {
		v.Set("http.port", port)
	}
	if logPath != "" {
		v.Set("zap.director", logPath)
	}

	var conf MigrateConf

	// 解析配置文件到结构体
	if err = v.Unmarshal(&conf); err != nil {
		fmt.Println("parse config file error: ", err)
	}

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
