/*
@Time : 2019-06-15 00:21
@Author : yangping
@File : po_convert
@Desc :
*/
package confpo

type BinaryConvert struct {
	// 算法替换字符串
	BinaryStr string `yaml:"binaryStr"`
	// 长链接过期时间
	LongUrlExpire int64 `yaml:"longUrlExpire"`
	// 短链接过期时间
	TinyUrlExpire int64 `yaml:"tinyUrlExpire"`
}
