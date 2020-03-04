/*
@Time : 2019-06-14 09:32
@Author : yangping
@File : index.go
@Desc :
*/
package tinyUrl

import (
	"github.com/gin-gonic/gin"
	"tinyUrl/common/middleware"
)

func InitTinyUrlRoute(router *gin.RouterGroup) {
	router.Use(middleware.TokenAuthMiddleware())
	tinyUrl := router.Group("/url")
	UrlRoute(tinyUrl)
}
