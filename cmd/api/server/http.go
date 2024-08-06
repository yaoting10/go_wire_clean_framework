package server

import (
	"github.com/gophero/goal/ginx"
	"github.com/gophero/goal/logx"
	"github.com/gophero/goal/mailx"
	"github.com/gophero/goal/redisx"
	"goboot/internal/config"
	"goboot/internal/handler"
	"goboot/internal/middleware"
	"goboot/internal/route"

	"github.com/gin-gonic/gin"
)

type App struct {
	Engine *gin.Engine
	Conf   config.Conf
}

func NewServerHTTP(
	_ mailx.Conf,
	logger *logx.Logger,
	conf config.Conf,
	redis redisx.Client,

	// api
	dm *handler.DemoApi,
	uh *handler.UserHandler,
	lh *handler.LoginHandler,
	ssh *handler.SysSettingHandler,
) *App {
	r := gin.New()
	// 不使用代理，https://pkg.go.dev/github.com/gin-gonic/gin#section-readme
	err := r.SetTrustedProxies(nil)
	if err != nil {
		panic(err)
	}

	// 加载 demo
	if gin.Mode() == gin.DebugMode {
		r.LoadHTMLGlob(conf.Asset().Page.Path + "/*")
	}

	middleware.InitSkip(conf)
	middleware.InitFixedTokenUrl(conf)
	middleware.InitI18N(conf)
	//middleware.InitXxlJob(r, conf, logger)

	registerMiddleWare(r, logger, conf, redis)
	route.RegisterRouters(r, dm, uh, lh, ssh)

	// No route group has permission
	noAuthRouter := r.Group("/")
	{
		noAuthRouter.GET("/", func(ctx *gin.Context) {
			logger.WithContext(ctx).Info("hello")
			ginx.HandleSuccess(ctx, map[string]interface{}{
				"say": "Hi!",
			})
		})
	}
	return &App{Engine: r, Conf: conf}
}

func registerMiddleWare(r *gin.Engine, logger *logx.Logger, conf config.Conf, redis redisx.Client) {
	r.Use(
		//middleware.CORSMiddleware(),
		//middleware.ResponseLogMiddleware(logger),
		//middleware.RequestLogMiddleware(logger),
		middleware.AuthMiddleware(logger, conf, redis),
		middleware.Recover(logger),
		middleware.GinZapLogger(logger),
		middleware.GinI18nLocalize(),
	)
}
