package config

type Http struct {
	Ip   string `yaml:"ip"`
	Port int    `yaml:"port"`
	Env  string `yaml:"env"`
}
