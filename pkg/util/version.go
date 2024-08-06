package util

import (
	"github.com/gophero/goal/conv"
	"strings"
)

func GetReleaseVersion(version string) int {
	ver := strings.Split(version, ".")
	big := ver[0]
	// middle := ver[1]
	// v := conv.StringToInt(big)*10 + conv.StringToInt(middle)
	v := conv.StrToInt(big)
	return v
}
