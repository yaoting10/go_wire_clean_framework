package config

type Job struct {
	ServerUrl    string `mapstructure:"server-url" yaml:"server-url"`
	Token        string `mapstructure:"token" yaml:"token"`
	ExecutorName string `mapstructure:"executor-name" yaml:"executor-name"`
}
