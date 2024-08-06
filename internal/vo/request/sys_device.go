package request

type SysDeviceVo struct {
	DeviceId          string `form:"deviceId"`    // 设备ID
	Platform          string `form:"platform"`    // 平台，mac 或 android
	Channel           string `form:"channel"`     // app渠道
	VersionInfo       string `form:"versionInfo"` // 版本信息
	UserId            uint   `form:"userId"`
	DeviceModel       string `form:"deviceModel"` // 设备型号
	DeviceBrand       string `form:"deviceBrand"` // 设备品牌
	DeviceType        string `form:"deviceType"`  // 设备类型
	DeviceOrientation string `form:"deviceOrientation"`
	DevicePixelRatio  string `form:"devicePixelRatio"`
	System            string `form:"system"`  // 操作系统及版本
	SysLang           string `form:"sysLang"` // 操作系统语言
}
