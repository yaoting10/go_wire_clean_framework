package middleware

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/gophero/goal/errorx"
	"github.com/gophero/goal/jsonx"
	"github.com/gophero/goal/logx"
	"github.com/gophero/goal/redisx"
	"goboot/internal/config"
	"goboot/internal/consts"
	"goboot/pkg/util"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gophero/goal/assert"
	"github.com/gophero/goal/ciphers"
	"github.com/gophero/goal/conv"
	"github.com/gophero/goal/stringx"
)

const (
	redisLockPrefix = "lck_user_"
	uuidKey         = "uuid"
	maxLen          = 400
	FixTokenUserId  = 0
)

func AuthMiddleware(logger *logx.Logger, conf config.Conf, redis redisx.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := c.Request
		url := req.URL.Path
		ip := util.GetRealIp(c)

		var shortUrl = stringx.CutString(maxLen, url)
		if checkSkipAuth(url) {
			err := req.ParseForm()
			if err != nil {
				logger.Errorf("request skip auth url %s error: %v", url, err)
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			//logger.Infof("Request - Skipped URL: %s, param: %v", shortUrl, req.Form)
			c.Next()
			return
		}

		token := util.GetToken(c)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"msg": fmt.Sprintf("invalid token, URL: %s, token: %s", shortUrl, token)})
			return
		}
		if token == conf.Server().AuthConf.FixToken && !checkFixedTokenAuth(url) {
			c.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"msg": "invalid token:" + token})
			return
		}
		keyNum, _ := redis.Exists(c, token).Result()
		if keyNum < 1 {
			c.AbortWithStatus(http.StatusExpectationFailed) // token过期请求响应417
			return
		}
		userId := util.GetUserId(c, redis)
		if userId < FixTokenUserId {
			c.AbortWithStatus(http.StatusExpectationFailed)
			return
		}

		// 防重复提交
		m := c.Request.Method
		if m != http.MethodGet { // 非get请求需要防重复提交
			// 如果是固定token的请求，需要按照ip来判断是否是重复提交
			var key string
			if userId == FixTokenUserId {
				key = redisLockPrefix + stringx.String(ip)
			} else {
				key = redisLockPrefix + stringx.String(userId)
			}
			boolCmd := redis.SetNX(context.Background(), key, 1, time.Second*30) // 30s过期
			if boolCmd.Val() {                                                   // 设置成功，则请求结束后删除
				defer func() { // 不能单独抽取方法，必须在这个方法中定义
					redis.Del(context.Background(), key) // 解锁
				}()
			} else { // 重复请求
				logger.Warnf("⚠️️ 发现重复请求：userId - %d, url: %s", userId, url)
				//_ = c.Error(errorx.FrequentReq) // 重复请求直接响应429
				//c.AbortWithStatus()
				c.AbortWithStatusJSON(http.StatusTooManyRequests, errorx.FrequentReq.Error())
				return
			}
		}

		// 解码请求并重复绑定
		//global.logger.Infof("Request  - UserId: %-5d, URL: %-20s, Token: %v", userId, shortUrl, token)
		reqParam, err := decryptAndRebind(c, getCryptoParam(c, logger), req, int64(userId), token, logger)
		if err != nil {
			logger.Errorf("decrypt request parameters error: %v", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// 非 200 状态码不加密响应
		blw := &encryptRespWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer, token: token, c: c}
		c.Writer = blw
		c.Next()
		body := blw.body.String()
		tmp := stringx.CutString(maxLen, body)
		if shouldShowDetailLog(c) {
			if c.Writer.Status() == http.StatusOK {
				logger.Infof("Request - UserId: %-8v, URL: %v, param: %s, Token: %v, ip: %v", userId, shortUrl, reqParam, token, ip)
				logger.Infof("Resp-%v, UserId: %-8v, URL: %v, body: %v ", c.Writer.Status(), userId, shortUrl, tmp)
			} else {
				logger.Warnf("Warn Request - UserId: %-8v, URL: %v, param: %s, Token: %v, ip: %v", userId, shortUrl, reqParam, token, ip)
				logger.Warnf("Warn Resp-%v, UserId: %-8v, URL: %v, body: %v ", c.Writer.Status(), userId, shortUrl, tmp)
			}
		}
	}
}

func getCryptoParam(c *gin.Context, logger *logx.Logger) string {
	var crypto []string
	switch c.Request.Method {
	case http.MethodGet:
		fallthrough
	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		fallthrough
	case http.MethodDelete:
		// request上的 Form 包含 url 请求参数和 body 参数，直接获取
		v := c.Request.Header.Get("Content-Type")
		d, _, err := mime.ParseMediaType(v)
		if err != nil || !(d == "multipart/form-data") {
			if err := c.Request.ParseForm(); err != nil {
				logger.Errorf("parse form request error: %v", err)
				return crypto[0]
			}
			crypto = c.Request.Form[consts.EncryptParamKey]
		} else {
			mf, err := c.MultipartForm()
			if err != nil {
				logger.Errorf("parse multipart request error: %v", err)
				return ""
			}
			crypto = mf.Value[consts.EncryptParamKey]
		}
	default:
		logger.Error("unsupported http method: " + c.Request.Method)
	}
	return crypto[0]
}

// decryptAndRebind 将加密请求参数解密并重新绑定到请求流上以便后续读取
func decryptAndRebind(ctx *gin.Context, crypto string, req *http.Request, userId int64, token string, logger *logx.Logger) (string, error) {
	assert.Require(crypto != "", "invalid request, need sign")
	crypto = strings.TrimSpace(crypto)
	// 只有一个加密字符串
	key := requestKey(token)
	iv := []byte(key)
	//fmt.Println("(", crypto, ")")
	// 随机iv，每次的密文不同
	encBytes, err := hex.DecodeString(crypto) // 十六进制解码
	if err != nil {
		return "", err
	}
	rawBytes, err := ciphers.AES.Decrypt(encBytes, []byte(key), ciphers.CBC, iv) // aes解密
	if err != nil {
		return "", err
	}

	delete(req.Form, consts.EncryptParamKey)

	var shortUrl = cutString(req.URL.String())

	assert.Require(userId >= 0, "invalid request, need userId")
	if shouldShowDetailLog(ctx) {
		logger.Infof("Request - UserId: %d, URL: %s, param: %s, Token: %v", userId, shortUrl, cutString(string(rawBytes)), token)
	}
	// 解析参数并重写到request上
	params := make(map[string]any)
	_, err = jsonx.Parse(rawBytes, &params, jsonx.UseNumber)
	if err != nil {
		return "", err
	}

	if _, ok := params["uuid"]; !ok {
		assert.Require(ok, "invalid request, need uuid")
	}
	for k, v := range params {
		if _, b := req.Form[k]; !b {
			if v != nil {
				req.Form[k] = []string{fmt.Sprintf("%v", v)}
			}
		}
	}
	req.Form[consts.LoginUserIdKey] = []string{conv.Int64ToStr(userId)}
	return stringx.CutString(maxLen, string(rawBytes)), nil
}

func cutString(s string) string {
	if len(s) > maxLen {
		return s[:maxLen] + "..."
	}
	return s
}

func requestKey(token string) string {
	return strings.Split(token, "-")[0][0:16]
}

func responseKey(token string) string {
	s := strings.Split(token, "-")
	return s[len(s)-1][0:16]
}

type encryptRespWriter struct {
	gin.ResponseWriter
	body  *bytes.Buffer
	token string
	c     *gin.Context
}

func (w encryptRespWriter) Write(b []byte) (int, error) {
	st := w.c.Writer.Status()
	// 只有200的才加密响应
	if len(b) > 0 && st == http.StatusOK {
		w.body.Write(b) // 原始响应数据
		key := responseKey(w.token)
		iv := []byte(key)                                                     // 随机iv，每次的密文不同
		encBytes, err := ciphers.AES.Encrypt(b, []byte(key), ciphers.CBC, iv) // 加密
		if err != nil {
			return 0, err
		}

		ret := hex.EncodeToString(encBytes) // 编码为十六进制字符串
		//logger.Debugf("Response Encryped: %v", ret)
		return w.ResponseWriter.WriteString(ret)
	} else {
		w.body.Write(b)
	}
	return w.ResponseWriter.Write(b)
}
