package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gophero/goal/ginx"
	"github.com/gophero/goal/logx"
	"github.com/gophero/goal/mailx"
	"github.com/gophero/goal/redisx"
	"goboot/internal/config"
	"goboot/internal/middleware"
	"goboot/internal/route"
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
) *App {
	r := gin.New()
	// 不使用代理，https://pkg.go.dev/github.com/gin-gonic/gin#section-readme
	err := r.SetTrustedProxies(nil)
	if err != nil {
		panic(err)
	}

	middleware.InitSkip(conf)
	middleware.InitFixedTokenUrl(conf)
	middleware.InitI18N(conf)
	//middleware.InitXxlJob(r, conf, logger)

	registerMiddleWare(r, logger, conf, redis)
	route.RegisterRouters(r)

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
