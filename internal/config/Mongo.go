package config

type Mongo struct {
	Host       string `mapstructure:"host" yaml:"host"`               // 连接地址
	Port       string `mapstructure:"port" yaml:"port"`               // 端口
	AuthSource string `mapstructure:"auth-source" yaml:"auth-source"` // 权限
	Username   string `mapstructure:"username" yaml:"username"`       // 用户名
	Password   string `mapstructure:"password" yaml:"password"`       // 密码
	Database   string `mapstructure:"database" yaml:"database"`       // 数据库名
}
