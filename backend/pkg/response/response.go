package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Code 响应码
type Code int

const (
	SUCCESS Code = 200
	ERROR   Code = 500

	// 认证相关 1000+
	INVALID_TOKEN     Code = 1001
	TOKEN_EXPIRED     Code = 1002
	TOKEN_MISSING     Code = 1003
	INVALID_PASSWORD  Code = 1004
	USER_NOT_FOUND    Code = 1005
	USER_DISABLED     Code = 1006
	PERMISSION_DENIED Code = 1007

	// 文章相关 2000+
	POST_NOT_FOUND Code = 2001

	// 评论相关 3000+
	COMMENT_NOT_FOUND Code = 3001
)

// Response 统一响应结构
type Response struct {
	Code    Code        `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Result 返回响应
func Result(c *gin.Context, code Code, data interface{}, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	Result(c, SUCCESS, data, "success")
}

// Error 错误响应
func Error(c *gin.Context, code Code, message string) {
	if message == "" {
		message = code.String()
	}
	Result(c, code, nil, message)
}

// String 返回码转字符串
func (c Code) String() string {
	switch c {
	case SUCCESS:
		return "success"
	case ERROR:
		return "error"
	case INVALID_TOKEN:
		return "invalid token"
	case TOKEN_EXPIRED:
		return "token expired"
	case TOKEN_MISSING:
		return "token missing"
	case INVALID_PASSWORD:
		return "invalid password"
	case USER_NOT_FOUND:
		return "user not found"
	case USER_DISABLED:
		return "user disabled"
	case PERMISSION_DENIED:
		return "permission denied"
	case POST_NOT_FOUND:
		return "post not found"
	case COMMENT_NOT_FOUND:
		return "comment not found"
	default:
		return "unknown error"
	}
}

// PageResult 分页响应
type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

func PageSuccess(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	Success(c, PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}
