package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gophero/goal/ciphers"
	"github.com/gophero/goal/logx"
	"github.com/gophero/goal/stringx"
	"github.com/gophero/goal/uuid"
	"go.uber.org/zap"
	"goboot/pkg/util"
	"strconv"
	"time"
)

func RequestLogMiddleware(logger *logx.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if shouldShowDetailLog(ctx) {
			// The configuration is initialized once per request
			trace := ciphers.MD5(uuid.UUID())
			logger.NewContext(ctx, zap.String("trace", trace))
			logger.NewContext(ctx, zap.String("request_method", ctx.Request.Method))
			headers, _ := json.Marshal(ctx.Request.Header)
			logger.NewContext(ctx, zap.String("request_headers", string(headers)))
			logger.NewContext(ctx, zap.String("request_url", ctx.Request.URL.String()))
			//if ctx.Request.Body != nil {
			//bodyBytes, _ := ctx.GetRawData()
			//ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 关键点
			logger.NewContext(ctx, zap.String("request_params", "...")) // 不打印加密的密文
			//}
			logger.WithContext(ctx).Info("Request")
		}
		ctx.Next()
	}
}

func ResponseLogMiddleware(logger *logx.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if shouldShowDetailLog(ctx) {
			blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
			ctx.Writer = blw
			startTime := time.Now()
			ctx.Next()
			duration := int(time.Since(startTime).Milliseconds())
			ctx.Header("X-Response-Time", strconv.Itoa(duration))
			logger.WithContext(ctx).Info("Response", zap.Any("response_body", blw.body.String()), zap.Any("time", fmt.Sprintf("%sms", strconv.Itoa(duration))))
		}
	}
}

func shouldShowDetailLog(ctx *gin.Context) bool {
	url := ctx.Request.URL.Path
	return !checkSkipAuth(url) && !stringx.StartsWith(url, "/user/report")
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// GinZapLogger 使用 zap 来接收gin框架的日志
func GinZapLogger(log *logx.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.RequestURI()
		// query := c.Request.URL.RawQuery

		c.Next()

		cost := time.Since(start).Microseconds()
		if shouldShowDetailLog(c) {
			log.Infof("%-4d %-5s %-18s %-30s %-4d | %-50s %s",
				c.Writer.Status(),
				c.Request.Method,
				// c.ClientIP(),
				util.GetRealIp(c),
				path,
				cost,
				c.Request.UserAgent(),
				c.Errors.ByType(gin.ErrorTypePrivate).String(),
			)
		}
	}
}
