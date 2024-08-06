package config

type Sign struct {
	Key string `mapstructure:"key" yaml:"key"`
}

type Auth struct {
	SkipUrl       string `mapstructure:"skip-url" yaml:"skip-url"`
	FixedTokenUrl string `mapstructure:"fixed-token-url" yaml:"fixed-token-url"`
	CaptchaKey    string `mapstructure:"captcha-key" yaml:"captcha-key"`
	FixToken      string `mapstructure:"fix-token" yaml:"fix-token"`
}
