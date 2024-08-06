package vo

import (
	"time"
)

const TOKEN = "token"

// Page 分页
type Page struct {
	PageNum  int `form:"pageNum" json:"pageNum"`
	PageSize int `form:"pageSize" json:"pageSize"`
}

type PageResp struct {
	Page
	Content       interface{} `form:"content" json:"content"`
	TotalElements int64       `form:"totalElements" json:"totalElements"`
}

type Model struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
