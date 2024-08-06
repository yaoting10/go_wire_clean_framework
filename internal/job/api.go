package job

import (
	"context"
	"github.com/gophero/goal/logx"
	"github.com/xxl-job/xxl-job-executor-go"
)

type Job interface {
	Name() string
	Logger() *logx.Logger
}

type NormalJob interface {
	Job
	Action(ctx context.Context, params ...any) (string, error)
}

type XxlJob interface {
	Name() string
	Action(ctx context.Context, param *xxl.RunReq) string
}
