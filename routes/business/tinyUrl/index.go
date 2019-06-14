/*
@Time : 2019-06-14 09:32
@Author : yangping
@File : index.go
@Desc :
*/
package tinyUrl

import "github.com/gin-gonic/gin"

func InitTinyUrlRoute(router *gin.RouterGroup) {
	tinyUrl := router.Group("/url")
	UrlRoute(tinyUrl)
}
