package response

import (
	"github.com/gin-gonic/gin"
)

// Response 구조체 정의
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success 성공 응답을 위한 헬퍼 함수
func Success(c *gin.Context, statusCode int, message string, data interface{}) Response {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// Error 에러 응답을 위한 헬퍼 함수
func Error(c *gin.Context, statusCode int, message string, err error) Response {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   err.Error(),
	})
	return Response{
		Success: false,
		Message: message,
		Error:   err.Error(),
	}
}
