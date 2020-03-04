/*
@date : 2020/03/03
@author : YaPi
@desc :
*/
package global_po

// Jwt 配置
type Jwt struct {
	// Jwt Secret
	JwtSecret string `yaml:"jwtSecret"`
	// Jwt 默认超时时间(单位 s)
	JwtExpireTime int    `yaml:"jwtExpireTime"`
	Issuer        string `yaml:"issuer"`
	Secret        []byte
}
