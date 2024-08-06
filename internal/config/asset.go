package config

type Asset struct {
	Page Page
	I18n I18n
}

type Page struct {
	Path string `mapstructure:"path" yaml:"path"`
}
