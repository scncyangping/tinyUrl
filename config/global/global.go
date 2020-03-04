/*
@date : 2020/03/03
@author : YaPi
@desc :
*/
package global

import (
	"tinyUrl/common/constants"
	"tinyUrl/config"
	"tinyUrl/config/global/global_po"
)

var JwtData *global_po.Jwt

func Init() {
	if JwtData == nil {
		JwtData = &global_po.Jwt{}
		if config.Base.Jwt.JwtSecret == constants.EmptyStr {
			JwtData.JwtSecret = constants.JwtSecret
		}
		if config.Base.Jwt.JwtExpireTime == constants.ZERO {
			JwtData.JwtExpireTime = constants.JwtExpireTime
		}

		if config.Base.Jwt.Issuer == constants.EmptyStr {
			JwtData.Issuer = constants.Issuer
		}

		JwtData = &global_po.Jwt{
			JwtExpireTime: config.Base.Jwt.JwtExpireTime,
			JwtSecret:     config.Base.Jwt.JwtSecret,
			Issuer:        config.Base.Jwt.Issuer,
			Secret:        []byte(config.Base.Jwt.JwtSecret),
		}

	}
}
