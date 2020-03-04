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
	"time"
	"tinyUrl/common/constants"
	"tinyUrl/common/http"
	"tinyUrl/common/util"
	"tinyUrl/config/log"
	"tinyUrl/domain/dto"
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
			errFlag  int
			session  = &dto.Session{}
		)
		token = c.Request.FormValue("API_TOKEN")
		if token == constants.EmptyStr {
			token = c.GetHeader("API_TOKEN")
		}
		if token == constants.EmptyStr {
			c.Abort()
			response.SetCode(http.RequestTokenNotFount)
			response.SetMsg(http.StatusText(http.RequestTokenNotFount))
			http.SendFailureRep(c, response)
			return

		}

		if claims, err := util.ParseToken(token); err != nil {
			errFlag = http.RequestCheckTokenError
		} else if time.Now().Unix() > claims.ExpiresAt {
			errFlag = http.RequestCheckTokenTimeOut
		} else {
			// 合法
			// 设置登录信息到token里面
			session.UserName = claims.Username
			session.Password = claims.Password
			c.Set("Session", session)
		}
		if errFlag > constants.ZERO {
			c.Abort()
			response.SetCode(errFlag)
			response.SetMsg(http.StatusText(errFlag))
			http.SendFailureRep(c, response)
			return
		}

		c.Next()
	}

}
