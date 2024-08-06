package task

import (
	"github.com/gophero/goal/errorx"
	"github.com/gophero/goal/logx"
	"github.com/gophero/goal/queue"
	"time"
)

func showDataLenInQueue(name string, logger *logx.Logger, queue queue.Queue, max int) {
	go func() {
		defer errorx.Recover(logger)
		for {
			logger.Warnf("task count：%-6d, name: %-30s", queue.Len(), name)
			time.Sleep(time.Second * 10)
		}
	}()
}
