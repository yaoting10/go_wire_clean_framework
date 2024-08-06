package job

import "github.com/gophero/goal/logx"

type Logger struct {
	Log *logx.Logger
}

func (l *Logger) Info(format string, a ...interface{}) {
	l.Log.Debugf("[JOB] - "+format, a...)
}

func (l *Logger) Error(format string, a ...interface{}) {
	l.Log.Errorf("[JOB] - "+format, a...)
}
