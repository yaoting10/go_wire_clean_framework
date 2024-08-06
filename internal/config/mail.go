package config

type Mail struct {
	Server        string `yaml:"server"`
	SupportDomain string `yaml:"supportDomain"`
}
