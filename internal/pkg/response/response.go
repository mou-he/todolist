package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 1. 定义统一的响应结构体
type Response struct {
	Code int         `json:"code"`           // 业务状态码（0=成功，非0=失败）
	Msg  string      `json:"msg"`            // 提示信息
	Data interface{} `json:"data,omitempty"` // 数据（omitempty 表示无数据时不返回该字段）
}

// 2. 辅助函数：成功返回
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

// 3. 辅助函数：失败返回（带自定义错误信息）
func Error(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{ // 注意：HTTP 状态码依然返回 200，通过业务 code 区分
		Code: -1,
		Msg:  msg,
		Data: nil,
	})
}

// 4. 辅助函数：失败返回（带自定义业务错误码，高级用法）
func ErrorWithCode(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

// 5. 参数校验失败时用的快捷方法
func ValidationError(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, Response{ // 此时返回 400 状态码
		Code: 400,
		Msg:  msg,
		Data: nil,
	})
}
