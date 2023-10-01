package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func UseAppRouter(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "感谢使用 vgo 后台!\nhttps://github.com/codelunaticer/vgo.go")
	})
}
