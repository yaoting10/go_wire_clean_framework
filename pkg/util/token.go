package util

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gophero/goal/ciphers"
	"github.com/gophero/goal/conv"
	"github.com/gophero/goal/redisx"
	"github.com/gophero/goal/uuid"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func GetToken(c *gin.Context) string {
	a := strings.TrimSpace(c.GetHeader(tokenKey))
	ss := strings.Split(a, " ")
	if ss[0] != tokenPrefix {
		fmt.Printf("❌ get token error, invalid prefix: %s\n", ss[0])
		return ""
	}
	return ss[1]
}

func ProduceToken(client redisx.Client, c context.Context, userId uint) string {
	return ProduceTokenExpire("", client, c, userId, 0) // no expire
}

func ProduceTokenExpire(prefix string, client redisx.Client, c context.Context, userId uint, expireInSec int) string {
	prefix += ciphers.MD5(conv.Int64ToStr(int64(userId)))
	middle := ciphers.MD5(uuid.UUID32())
	suffix := ciphers.MD5(conv.Int64ToStr(time.Now().UnixMilli())) // 时间戳

	delKeys(client, c, prefix)
	token := prefix + "-" + middle + "-" + suffix
	client.Set(c, token, userId, time.Second*time.Duration(expireInSec))
	return token
}

func delKeys(client redisx.Client, ctx context.Context, keyPrefix string) bool {
	var doDelKey = func(c redisx.Client) {
		keys := client.Keys(ctx, keyPrefix+"*").Val()
		for _, k := range keys {
			client.Del(ctx, k)
		}
	}
	if cc, ok := client.(*redis.ClusterClient); ok {
		if err := cc.ForEachMaster(ctx, func(ctx context.Context, c *redis.Client) error {
			doDelKey(c)
			return nil
		}); err != nil {
			return false
		}
	} else {
		doDelKey(client)
	}
	return true
}
