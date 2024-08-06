package middleware

import (
	"github.com/gophero/goal/logx"
	"goboot/internal/config"
	"goboot/internal/job"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xxl-job/xxl-job-executor-go"
)

func InitXxlJob(r *gin.Engine, conf config.Conf, logger *logx.Logger, jobs ...job.XxlJob) xxl.Executor {
	exec := xxl.NewExecutor(
		xxl.ServerAddr(conf.Job().ServerUrl),
		xxl.AccessToken(conf.Job().Token),                                   // 请求令牌，与job服务端配置的令牌保持一致
		xxl.ExecutorIp(conf.GetViper().GetString("http.ip")),                // 执行器ip地址，即应用ip
		xxl.ExecutorPort(strconv.Itoa(conf.GetViper().GetInt("http.port"))), // 执行器端口，与应用端口保持一致
		xxl.RegistryKey(conf.Job().ExecutorName),                            // 执行器名称
		xxl.SetLogger(&job.Logger{Log: logger}),                             // 自定义日志
	)
	exec.Init()

	// 设置日志查看handler
	exec.LogHandler(func(req *xxl.LogReq) *xxl.LogRes {
		return &xxl.LogRes{Code: xxl.SuccessCode, Msg: "", Content: xxl.LogResContent{
			FromLineNum: req.FromLineNum,
			ToLineNum:   2,
			LogContent:  "这个是自定义日志handler",
			IsEnd:       true,
		}}
	})

	XxlJobMux(r, exec) // 集成gin, 注册任务 endpoint，不需要在调用 exec.Run() 了

	// 注册任务
	for _, j := range jobs {
		exec.RegTask(j.Name(), j.Action)
	}
	return exec
}

func XxlJobMux(e *gin.Engine, exec xxl.Executor) {
	// 注册的gin的路由
	// 暴露在外网，可能被恶意攻击，受accessToken机制保护
	e.POST("run", gin.WrapF(exec.RunTask))
	e.POST("kill", gin.WrapF(exec.KillTask))
	e.POST("log", gin.WrapF(exec.TaskLog))
	e.POST("beat", gin.WrapF(exec.Beat))
	e.POST("idleBeat", gin.WrapF(exec.IdleBeat))
}
