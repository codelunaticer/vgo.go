package httpLogic

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
	tipConst "vgoadmin-go/constant/tip"
	validatorHelper "vgoadmin-go/helper/validator"
)

type _response struct {
	Data    interface{} `json:"data"`
	Message interface{} `json:"message"`
	Code    int         `json:"code"`
}

// ErrorResponse 响应错误状态对象
func ErrorResponse(c *gin.Context, status int, message interface{}) {
	c.JSON(http.StatusBadRequest, _response{
		Message: message,
		Code:    status,
	})
}

// BadErrorResponse 参数验证失败响应对象
func BadErrorResponse(c *gin.Context, message interface{}) {
	c.JSON(http.StatusUnprocessableEntity, _response{
		Message: message,
		Code:    422,
	})
}

// NoAuthResponse 中间件拦截响应
func NoAuthResponse(c *gin.Context, status int, message interface{}) {
	c.JSON(http.StatusUnauthorized, _response{
		Message: message,
		Code:    status,
	})
}

func OKResponse(c *gin.Context, data interface{}, message interface{}) {
	if message == nil {
		message = tipConst.ResOK
	}
	c.JSON(http.StatusOK, _response{
		Data:    data,
		Message: message,
		Code:    200,
	})
}

// GetBindErrorTranslate 获取bind验证错误信息的中文翻译
func GetBindErrorTranslate(err error) string {
	errs, ok := err.(validator.ValidationErrors)
	if ok {
		// 对validationErrors类型错误进行翻译并合并成字符串
		var errorMsgs []string
		for _, e := range errs {
			errorMsgs = append(errorMsgs, e.Translate(validatorHelper.GetTrans()))
		}
		return strings.Join(errorMsgs, ", ") // 以逗号分隔错误消息
	}
	// 其他类型错误直接返回
	return err.Error()
}

// ExeErrorResponse 程序错误响应
func ExeErrorResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Internal Server Error"
	}
	// 返回自定义错误响应
	c.JSON(http.StatusInternalServerError, _response{
		Message: message,
		Code:    500,
	})
}
