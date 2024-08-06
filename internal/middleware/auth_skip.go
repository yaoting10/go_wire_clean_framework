package middleware

import (
	"goboot/internal/config"
	"regexp"
	"strings"
)

// 跳过认证的url

const (
	blurString       = "*"
	doubleBlurString = "**"
)

// skipUrlRegex 跳过认证的正则，从配置的 conf.Auth 的 skipUrl 编译
var skipUrlRegex []*regexp.Regexp
var fixedTokenUrlRegex []*regexp.Regexp

// InitSkip 初始化时编译 skipUrlRegex
func InitSkip(conf config.Conf) {
	skipUrl := conf.Server().AuthConf.SkipUrl
	if skipUrl != "" {
		skipUrls := strings.Split(skipUrl, ",")
		for _, su := range skipUrls {
			su = formatRegexString(su)
			r, _ := regexp.Compile(su)
			skipUrlRegex = append(skipUrlRegex, r)
		}
	}
}

// InitFixedTokenUrl 初始化时编译 skipUrlRegex
func InitFixedTokenUrl(conf config.Conf) {
	fixedTokenUrl := conf.Server().AuthConf.FixedTokenUrl
	if fixedTokenUrl != "" {
		skipUrls := strings.Split(fixedTokenUrl, ",")
		for _, su := range skipUrls {
			su = formatRegexString(su)
			r, _ := regexp.Compile(su)
			fixedTokenUrlRegex = append(fixedTokenUrlRegex, r)
		}
	}
}

// formatRegexString 将配置的 skipUrl 转换为正则形式
func formatRegexString(s string) string {
	s = "^" + s
	if i := strings.Index(s, doubleBlurString); i > 0 {
		// /**/
		s = strings.ReplaceAll(s, "/"+doubleBlurString+"/", "&&&")
		// /**
		s = strings.ReplaceAll(s, "/"+doubleBlurString, "&&")
	}
	if i := strings.Index(s, blurString); i > 0 {
		s = strings.ReplaceAll(s, blurString, "&")
	}
	s = strings.ReplaceAll(s, "&&&", "/(\\w|/?)*")
	s = strings.ReplaceAll(s, "&&", "/(\\w|/)*")
	s = strings.ReplaceAll(s, "&", "(\\w)+")
	return s + "$"
}

// checkSkipAuth 根据指定的 url 检查是否跳过认证，跳过则返回 true，否则返回 false
func checkSkipAuth(url string) bool {
	for _, regex := range skipUrlRegex {
		if regex.Match([]byte(url)) {
			return true
		}
	}
	return false
}

// checkFixedTokenAuth 根据固定的token和请求url 检查是否跳过认证，跳过则返回 true，否则返回 false
func checkFixedTokenAuth(url string) bool {
	for _, regex := range fixedTokenUrlRegex {
		if regex.Match([]byte(url)) {
			return true
		}
	}
	return false
}
