package repo

// PageResult 分页结果对象
type PageResult[T any] struct {
	Data []T // 当前页数据
}

// ScrollPageResult 滚动分页
type ScrollPageResult[T any] struct {
	PageResult[T]
	HasNext bool
	HasPrev bool
}

// Pager 分页接口
// type Pager[T any] interface {
// 	Paging(request PageRequest) (PageResult[T], error)
// }

// PageRequest 分页请求
type PageRequest struct {
	Index int // 当前页数
	Size  int // 每一页条数
}

func NewPageRequest(index, size int) *PageRequest {
	return &PageRequest{Index: index, Size: size}
}

func (pr *PageRequest) Start() int {
	return pr.Index * pr.Size
}
