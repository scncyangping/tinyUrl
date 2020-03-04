/*
@Time : 2019-06-14 16:18
@Author : yangping
@File : tinyInfo
@Desc :
*/
package entity

type TinyInfo struct {
	Id         string `form:"id" json:"id" xml:"id" bson:"_id"`
	UserName   string `form:"userName" json:"id" userName:"id" bson:"userName"`
	LongUrl    string `form:"longUrl" json:"longUrl" xml:"longUrl"`
	TinyUrl    string `form:"tinyUrl" json:"tinyUrl" xml:"tinyUrl"`
	Count      int    `form:"count" json:"count" xml:"count"`
	Type       string `form:"type" json:"type" xml:"type" bson:"type"`
	CreateTime int64  `form:"createTime" json:"createTime" xml:"createTime"`
	ExpireTime int    `form:"expireTime" json:"expireTime" xml:"expireTime"`
}
