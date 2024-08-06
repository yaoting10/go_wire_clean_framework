package config

import "github.com/gophero/goal/twitter"

type Twitter struct {
	twitter.Setting `mapstructure:",squash"`
	RedirectUrl     string `mapstructure:"redirectUrl" yaml:"redirectUrl"`
}
