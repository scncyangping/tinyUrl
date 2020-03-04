/*
@Time : 2019-06-14 10:11
@Author : yangping
@File : tinyUrl
@Desc :
*/
package tinyUrl

import (
	"github.com/gin-gonic/gin"
	"tinyUrl/handler/business"
)

func UrlRoute(router *gin.RouterGroup) {
	router.POST("/", business.UrlTransform)

	router.POST("/custom", business.UrlTransformCustom)

	router.GET("/info", business.UrlBaseInfo)
}
