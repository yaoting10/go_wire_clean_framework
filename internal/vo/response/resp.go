package response

import (
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Resp struct {
	*gin.Context
}

type Map map[string]any

func Ctx(c *gin.Context) *Resp {
	return &Resp{c}
}

func (r *Resp) Ok() {
	r.Status(http.StatusOK)
}

func (r *Resp) OkJson(data any) {
	r.JSON(http.StatusOK, data)
}

func (r *Resp) Okf(s string, vs ...any) {
	r.JSON(http.StatusOK, fmt.Sprintf(s, vs...))
}

func (r *Resp) OkI18n(key string, vs ...any) {
	r.JSON(http.StatusOK, fmt.Sprintf(i18n.MustGetMessage(key), vs...))
}

func (r *Resp) ServerErr() {
	r.Status(http.StatusInternalServerError)
}

func (r *Resp) ServerErrf(fmt string, vs ...any) {
	r.String(http.StatusInternalServerError, fmt, vs...)
}

func (r *Resp) ServerI18n(key string, vs ...any) {
	r.String(http.StatusInternalServerError, i18n.MustGetMessage(key), vs...)
}

func (r *Resp) NotFound() {
	r.Status(http.StatusNotFound)
}

func (r *Resp) NotFoundf(fmt string, vs ...any) {
	r.String(http.StatusNotFound, fmt, vs...)
}

func (r *Resp) NotFoundI18n(key string, vs ...any) {
	r.String(http.StatusNotFound, i18n.MustGetMessage(key), vs...)
}

func (r *Resp) BadReq() {
	r.Status(http.StatusBadRequest)
}

func (r *Resp) BadReqf(fmt string, vs ...any) {
	r.String(http.StatusBadRequest, fmt, vs...)
}

func (r *Resp) BadReqI18n(key string, vs ...any) {
	r.String(http.StatusBadRequest, i18n.MustGetMessage(key), vs...)
}

func (r *Resp) NoAuth() {
	r.Status(http.StatusUnauthorized)
}

func (r *Resp) NoAuthf(fmt string, vs ...any) {
	r.String(http.StatusUnauthorized, fmt, vs...)
}

func (r *Resp) NoAuthI18n(key string, vs ...any) {
	r.String(http.StatusUnauthorized, i18n.MustGetMessage(key), vs...)
}
