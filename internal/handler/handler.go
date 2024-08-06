package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gophero/goal/logx"
)

type Handler struct {
	L *logx.Logger
}

func NewHandler(logger *logx.Logger) *Handler {
	return &Handler{
		L: logger,
	}
}

func (h *Handler) Group(path string, e *gin.Engine, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return e.Group(path, handlers...)
}
