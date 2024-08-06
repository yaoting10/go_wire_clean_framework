package middleware

//
//import (
//	"github.com/gin-gonic/gin"
//	"global"
//	"time"
//	"video-api/g"
//	myi18n "video-api/i18n"
//	"video-api/job"
//)
//
//type RegisterJob interface {
//	RegisterRouter(r *gin.Engine)
//}
//
//var routers []RegisterJob
//
//// RegisterRouters 注册路由表
//func RegisterRouters(r RegisterJob) {
//	routers = append(routers, r)
//}
//
//func InitRouter() *gin.Engine {
//	r := gin.New()
//	// 不使用代理，https://pkg.go.dev/github.com/gin-gonic/gin#section-readme
//	err := r.SetTrustedProxies(nil)
//	if err != nil {
//		panic(err)
//	}
//
//	// 加载 demo 页面
//	if gin.Mode() == gin.DebugMode {
//		r.LoadHTMLGlob("pages/*")
//	}
//
//	initI18N(r)           // 初始化 i18n
//	InitSkip()            // 初始化认证跳过器
//	InitFixedTokenUrl()   // 初始化固定token跳过
//	registerMiddleWare(r) // 注册中间件
//	registerRouters(r)    // 注册路由
//	job.InitXxlJob(r)           // 初始化调度任务
//
//	return r
//}
//
//func initI18N(r *gin.Engine) {
//	myi18n.InitXxlJob(g.Conf.I18n.Path)
//}
//
//func registerRouters(e *gin.Engine) {
//	for _, r := range routers {
//		r.RegisterRouter(e)
//	}
//}
//
//func registerMiddleWare(r *gin.Engine) {
//	r.Use(AuthHandler())            // token机制
//	r.Use(GinLogger())              // Gin日志
//	r.Use(Recover)                  // 异常恢复
//	r.Use(myi18n.GinI18nLocalize()) // 国际化
//}
//
//// GinLogger 使用 zap 来接收gin框架的日志
//func GinLogger() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		start := time.Now()
//		path := c.Request.URL.RequestURI()
//		// query := c.Request.URL.RawQuery
//
//		c.Next()
//
//		cost := time.Since(start).Microseconds()
//
//		global.logger.Infof("%-4d %-5s %-18s %-30s %-4d | %-50s %s",
//			c.Writer.Status(),
//			c.Request.Method,
//			// c.ClientIP(),
//			getRealIp(c),
//			path,
//			cost,
//			c.Request.UserAgent(),
//			c.Errors.ByType(gin.ErrorTypePrivate).String(),
//		)
//	}
//}
