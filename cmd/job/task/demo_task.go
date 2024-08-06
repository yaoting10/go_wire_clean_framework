package task

import (
	"context"
	"fmt"
	"github.com/gophero/goal/queue"
	"github.com/gophero/goal/redisx"
	"github.com/xxl-job/xxl-job-executor-go"
	"goboot/internal/service"
	"time"
)

const (
	groupCanDismissCheckKey = "queue:check_group_dismiss"
	demoTaskQueueThreshold  = 1000 // 队列中的元素个数阈值
)

type DemoTask struct {
	*service.Service
	rc    redisx.Client
	queue queue.Queue
}

func NewDemoTask(srv *service.Service, rc redisx.Client) *DemoTask {
	job := &DemoTask{Service: srv, rc: rc,
		queue: redisx.NewBoundedUniqueQueue(rc, groupCanDismissCheckKey, demoTaskQueueThreshold),
	}
	return job
}

func (job *DemoTask) Name() string {
	return "job.demo_task"
}

func (job *DemoTask) Action(ctx context.Context, param *xxl.RunReq) string {
	time.Sleep(30 * time.Second)
	fmt.Println("demo task success")
	return "ok"
}
