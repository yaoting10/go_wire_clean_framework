package request

type GuideTask struct {
	Request
	TaskId int `form:"taskId"`
}
