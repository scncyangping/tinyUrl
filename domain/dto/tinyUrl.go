/*
@Time : 2019-06-14 17:10
@Author : yangping
@File : tinyUrl
@Desc :
*/
package dto

type TinyDto struct {
	LongUrl    string `form:"longUrl" json:"longUrl" xml:"longUrl"`
	TinyUrl    string `form:"tinyUrl" json:"tinyUrl" xml:"tinyUrl"`
	ExpireTime int    `form:"expireTime" json:"expireTime" xml:"expireTime"`
}
