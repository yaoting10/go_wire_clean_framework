package config

type Google struct {
	ClientId     string `mapstructure:"client-id" yaml:"client-id"`
	ClientSecret string `mapstructure:"client-secret" yaml:"client-secret"`
	RedirectUrl  string `mapstructure:"redirect-url" yaml:"redirect-url"`
	CaptchaKey   string `mapstructure:"captcha-key" yaml:"captcha-key"`
}
