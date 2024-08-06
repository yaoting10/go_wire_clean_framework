package consts

const (
	_                = iota
	AppChannelApk    //官网apk包
	AppChannelAab    //谷歌商店aab包
	AppChannelIosTf  //ios的TF包
	AppChannelIosApp //ios商店包
)

const (
	AppChannelKey          = "AppChannel"
	AppChanAndroid         = "android"
	AppChanGooglePlayStore = "google"
	AppChanIOSTestFlight   = "testflight"
	AppChanIOSAppStore     = "ios"
)

const (
	DefaultLanguageCode = "en"
)

func GetChanName(ch int) string {
	switch ch {
	case AppChannelApk:
		return AppChanAndroid
	case AppChannelAab:
		return AppChanGooglePlayStore
	case AppChannelIosTf:
		return AppChanIOSTestFlight
	case AppChannelIosApp:
		return AppChanIOSAppStore
	}
	return ""
}
