/*
@Time : 2019/4/30 11:24 AM
@Author : yangping
@File : middleware
@Desc :
*/
package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"tinyUrl/common/constants"
	"tinyUrl/common/http"
	"tinyUrl/common/util"
	"tinyUrl/config/log"
)

/*
 * date : 2019/4/30
 * author : yangping
 * desc : 增加统一日志
 */
func VisitLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			actionId string
		)
		actionId = util.UUID()
		data, err := ctx.GetRawData()
		if err != nil {
			log.GetLogger().Errorf("Visit Param Init Error %v", err)
		}
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))

		log.GetLogger().Infof(
			"%s request %s [%s] from [%s]: %s",
			ctx.Request.Method,
			ctx.Request.RequestURI,
			actionId, ctx.ClientIP(),
			string(data))

		ctx.Set("ActionId", actionId)
		ctx.Next()
	}
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			token    string
			response = http.Instance()
		)
		token = c.Request.FormValue("API_TOKEN")
		if token == constants.EmptyStr {
			token = c.GetHeader("API_TOKEN")
		}
		if token == constants.EmptyStr {
			c.Abort()
			response.SetCode(http.RequestCheckTokenError)
			response.SetMsg(http.StatusText(http.RequestCheckTokenError))
			http.SendFailureRep(c, response)
			return

		}

		// TEST
		if token != "API_TOKEN" {
			c.Abort()
			response.SetCode(http.RequestCheckTokenError)
			response.SetMsg(http.StatusText(http.RequestCheckTokenError))
			http.SendFailureRep(c, response)
			return
		}
		c.Next()
	}

}
