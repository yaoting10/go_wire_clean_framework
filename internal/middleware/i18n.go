package middleware

import (
	"fmt"
	"goboot/internal/config"
	"os"
	"strings"

	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

var bundleConfig *i18n.BundleCfg
var supportLangs map[string]lang
var defaultLanguage = language.English

type lang struct {
	lang language.Tag
	code string
	name string
}

func InitI18N(conf config.Conf) {
	var (
		defaultUnmarshalFunc  = yaml.Unmarshal
		defaultAcceptLanguage []language.Tag
		defaultLoader         = i18n.LoaderFunc(os.ReadFile)
	)
	supportLangs = make(map[string]lang, 10)
	supportLangs["en"] = lang{lang: language.English, code: "en", name: "英语"}
	//supportLangs["es"] = lang{lang: language.Spanish, code: "es", name: "西班牙语"}
	//supportLangs["fr"] = lang{lang: language.French, code: "fr", name: "法语"}
	//supportLangs["hi"] = lang{lang: language.Hindi, code: "hi", name: "印地语"}
	//supportLangs["ru"] = lang{lang: language.Russian, code: "ru", name: "俄语"}
	//supportLangs["vi"] = lang{lang: language.Vietnamese, code: "vi", name: "越南语"}
	//supportLangs["ja"] = lang{lang: language.Japanese, code: "ja", name: "日语"}
	supportLangs["zh"] = lang{lang: language.SimplifiedChinese, code: "zh", name: "简体中文"} // 语言文件名必须是: zh-Hans.yaml

	for _, v := range supportLangs {
		fmt.Printf("\tlang : %-6s, code: %-6s, name: %s\n", v.lang.String(), v.code, v.name)
		defaultAcceptLanguage = append(defaultAcceptLanguage, v.lang)
	}
	// Bundle配置
	bundleConfig = &i18n.BundleCfg{
		RootPath:         conf.I18n().Path,
		AcceptLanguage:   defaultAcceptLanguage,
		FormatBundleFile: "yaml",
		DefaultLanguage:  defaultLanguage,
		UnmarshalFunc:    defaultUnmarshalFunc,
		Loader:           defaultLoader,
	}
}

func getLangTag(code string) string {
	if code != "" {
		code = strings.Split(code, "-")[0]
		code = strings.Split(code, "_")[0]
		code = strings.ToLower(code)
		if mp, ok := supportLangs[code]; ok {
			if mp.code != "" {
				return mp.lang.String()
			}
		}
	}
	return defaultLanguage.String()
}

func SupportLang(code string) bool {
	code = strings.Split(code, "-")[0]
	code = strings.Split(code, "_")[0]
	_, ok := supportLangs[code]
	return ok
}

// GinI18nLocalize gin i18n中间件
func GinI18nLocalize() gin.HandlerFunc {
	return i18n.Localize( /* 使用配置的Bundle */
		i18n.WithBundle(bundleConfig),
		// 请求时处理语言
		i18n.WithGetLngHandle(func(context *gin.Context, defaultLng string) string {
			// context.GetHeader 抛出空指针异常，捕获异常，使用默认语言
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("i18n error: %v\n", err)
				}
			}()
			alang := context.GetHeader("Accept-Language")
			alang = getLangTag(alang)
			//req := request.Bind(context, request.Request{})
			//fmt.Printf("userId: %d, url: %s, code: %s, lang: %s\n", req.LoginUserId,
			//	context.Request.URL.Path,
			//	context.GetHeader("Accept-Language"), alang)
			//if alang == "" {
			return defaultLng
			//}
			//return alang
		}),
	)
}
