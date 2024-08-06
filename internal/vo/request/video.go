package request

import (
	"goboot/internal/vo"
	"time"
)

type VideoRequest struct {
	Request
	vo.Page
	Language string
	VideoId  []string `form:"videoId" json:"videoId"`
	Type     int8     `form:"type" json:"type"` // 0查询所有，1查询关注视频列表，2查单个用户的视频列表
	UserId   uint     `form:"userId" json:"userId"`
	UserIds  []uint   `form:"userIds" json:"userIds"`
	Keyword  string   `form:"keyword" json:"keyword"`
}

type DelVideoReq struct {
	Request
	Id uint `json:"id" form:"id"`
}

type VideoUploadReq struct {
	Request
	Info string `json:"info" form:"info"`
}

type VideoInfo struct {
	ID          uint       `json:"id" form:"id"`
	Title       string     `json:"title" form:"title"`         // 标题
	Desc        string     `json:"desc" form:"desc"`           // 描述
	FileName    string     `json:"fileName" form:"fileName"`   // 全路径文件名
	CateId      uint       `json:"cateId" form:"cateId"`       // 分类主键id
	AliCateId   int        `json:"aliCateId" form:"aliCateId"` // 阿里云分类id
	Cover       string     `json:"cover" form:"cover"`         // 封面url
	CoverName   string     `json:"coverName" form:"coverName"` // 封面文件名称
	Tags        string     `json:"tags" form:"tags"`           // 标签
	VideoId     string     `json:"videoId" form:"videoId"`     // 视频id
	Lang        string     `json:"lang" form:"lang"`           // 语言
	UserId      uint       `json:"userId" form:"userId"`       // 用户id
	Url         string     `json:"url" form:"url"`             // 播放url
	Size        int32      `json:"size" form:"size"`
	Type        string     `json:"type" form:"type"`
	WorkflowId  string     `json:"workflowId" form:"workflowId"` // workflowId
	Weight      int        `json:"weight" form:"weight"`
	Width       int64      `json:"width" form:"width"`
	Height      int64      `json:"height" form:"height"`
	During      float64    `json:"during" form:"during"`
	Source      int        `json:"source" form:"source"`
	SourceUrl   string     `json:"sourceUrl" form:"sourceUrl"`
	ExpiredAt   *time.Time `json:"expiredAt" form:"expiredAt"`
	SourceFmtId string     `json:"sourceFmtId" form:"sourceFmtId"`
	Expired     int64      `json:"expired" form:"expired"`
}
