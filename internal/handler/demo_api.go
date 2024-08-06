package handler

import (
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/gophero/goal/random"
	"goboot/internal/service"
	"goboot/internal/vo/response"
	"math/rand"
	"net/http"
	"time"
)

type DemoApi struct {
	srv *service.Service
	*Handler
}

func NewDemoApi(h *Handler, srv *service.Service) *DemoApi {
	return &DemoApi{Handler: h, srv: srv}
}

func (api *DemoApi) RegisterRouter(r *gin.Engine) {
	//var g = api.Group("/demo", r)
	//{
	//g.GET("/hello", api.Hello)
	//g.POST("/test_fix_data", h.testFixData) // 转出nft

	//}
}

func (api *DemoApi) Hello(c *gin.Context) {
	i := rand.Int()
	api.L.Info("random int: ", i)
	c.String(http.StatusOK, "hello, %d, %s, text: %s", i,
		api.srv.C.GetViper().Get("http.port"),
		i18n.MustGetMessage("hello"),
	)
}

func (api *DemoApi) Panic(c *gin.Context) {
	i := rand.Int()
	panic(fmt.Errorf("make an error: %d, %s", i, api.srv.C.GetViper().Get("http.port")))
}

func (api *DemoApi) I18n(c *gin.Context) {
	// 模拟线上响应时间
	time.Sleep(time.Millisecond * time.Duration(200+random.Int(500)))
	response.Ctx(c).OkJson(gin.H{
		"lang": c.GetHeader("Accept-Language"),
		"text": i18n.MustGetMessage("hello"),
	})
}
