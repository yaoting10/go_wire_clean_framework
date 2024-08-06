package route

import (
	"github.com/gin-gonic/gin"
)

type Register interface {
	RegisterRouter(r *gin.Engine)
}

//var routers []Register

// RegisterRouter 注册路由表
//func RegisterRouter(r Register) {
//	routers = append(routers, r)
//}

func RegisterRouters(e *gin.Engine, rs ...Register) {
	//routers = append(routers, rs...)
	for _, r := range rs {
		r.RegisterRouter(e)
	}
}

//func registerRouters(e *gin.Engine) {
//	for _, r := range routers {
//		r.RegisterRouter(e)
//	}
//}
