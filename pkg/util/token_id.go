package util

import (
	"fmt"
	"github.com/gophero/goal/conv"
	"math"
	"math/rand"
	"strings"
	"time"
)

func GenerateTokenId(userId uint, selfInc int64) string {
	idLen := 3
	incLen := 2
	randLen := 3

	rand.Seed(time.Now().UnixNano())

	var builder strings.Builder

	idStr := conv.Int64ToStr(int64(userId))
	// idp := idStr[:1]
	// idStr = idp + conv.Int64ToString(int64(rand.Intn(10))) + conv.Int64ToString(int64(rand.Intn(10)))
	switch {
	case userId > 999:
		idp := idStr[:1]
		idStr = idp + idStr[len(idStr)-2:] // 第一位加后两位
	default:
		// id补零
		for {
			if len(idStr) >= idLen {
				break
			}
			idStr += "0"
		}
	}
	builder.WriteString(idStr)

	incStr := conv.Int64ToStr(selfInc)
	incStr = ProcessFillZero(incStr, incLen)
	builder.WriteString(incStr)

	randStr := conv.IntToStr(int(time.Now().UnixNano() / 1000 % 1000)) // 使用纳秒后3位
	randStr = ProcessFillZero(randStr, randLen)
	builder.WriteString(randStr)
	return builder.String()
}

func GenerateCert(userId uint) string {
	idLen := 5
	randLen := 3

	rand.Seed(time.Now().UnixNano())

	var builder strings.Builder

	idStr := conv.Int64ToStr(int64(userId))
	// idp := idStr[:1]
	// idStr = idp + conv.Int64ToString(int64(rand.Intn(10))) + conv.Int64ToString(int64(rand.Intn(10)))
	switch {
	case userId > 99999:
		idp := idStr[:1]
		idStr = idp + idStr[len(idStr)-4:] // 第一位位加后四位
	default:
		// id补零
		for {
			if len(idStr) >= idLen {
				break
			}
			idStr += "0"
		}
	}
	builder.WriteString(idStr)

	randStr := conv.IntToStr(int(time.Now().UnixNano() / 1000 % 1000)) // 使用纳秒后3位
	randStr = ProcessFillZero(randStr, randLen)
	builder.WriteString(randStr)
	return builder.String()
}

func RandomInt(begin int, end int) int {
	rand.Seed(time.Now().UnixNano())
	return begin + rand.Intn(end-begin)
}

func ZeroLeftFill(str string, count int) string {
	return fmt.Sprintf("%0*s", count, str)
}

func ProcessFillZero(str string, length int) string {
	if len(str) > length {
		str = str[len(str)-length:]
	} else if len(str) < length {
		str = ZeroLeftFill(str, length)
	}
	return str
}

func GeoLatLonDistance(lon1, lat1, lon2, lat2 float64) float64 {
	if lon1 == lon2 && lat1 == lat2 {
		return 0
	}
	radius := 6371000.0 // 6378137.0
	rad := math.Pi / 180.0
	lat1 = lat1 * rad
	lon1 = lon1 * rad
	lat2 = lat2 * rad
	lon2 = lon2 * rad
	theta := lon2 - lon1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	return dist * radius
}

func GenerateTxNum(txType int) string {
	s := fmt.Sprintf("%s%v%v%s", "T", txType, time.Now().UnixMilli(), RandStr(6))
	return s
}

func RandStr(n int) string {
	const numLetterString = "0123456789"
	str := []byte(numLetterString)
	res := ""
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		res += fmt.Sprintf("%c", str[rand.Intn(strings.Count(numLetterString, "")-1)])
	}
	return res
}

func RandomOne(sources []string) string {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(sources))
	return sources[i]
}
