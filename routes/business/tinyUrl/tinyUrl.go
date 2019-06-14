/*
@Time : 2019-06-14 10:11
@Author : yangping
@File : tinyUrl
@Desc :
*/
package tinyUrl

import (
	"github.com/gin-gonic/gin"
	"tinyUrl/handler/business/tinyHandler"
)

func UrlRoute(router *gin.RouterGroup) {
	router.GET("/", tinyHandler.UrlTransform)
	router.GET("/custom", tinyHandler.UrlTransformCustom)
	router.GET("/go", tinyHandler.Redirect4TinyUrl)
}
