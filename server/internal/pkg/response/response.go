package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Body struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Body{Code: 0, Msg: "ok", Data: data})
}

func Fail(c *gin.Context, httpStatus int, code int, msg string) {
	c.JSON(httpStatus, Body{Code: code, Msg: msg})
}

func BadRequest(c *gin.Context, msg string) {
	Fail(c, http.StatusBadRequest, 400, msg)
}

func Unauthorized(c *gin.Context, msg string) {
	Fail(c, http.StatusUnauthorized, 401, msg)
}

func Forbidden(c *gin.Context, msg string) {
	Fail(c, http.StatusForbidden, 403, msg)
}

func NotFound(c *gin.Context, msg string) {
	Fail(c, http.StatusNotFound, 404, msg)
}

func ServerError(c *gin.Context, msg string) {
	Fail(c, http.StatusInternalServerError, 500, msg)
}
