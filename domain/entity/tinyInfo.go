/*
@Time : 2019-06-14 16:18
@Author : yangping
@File : tinyInfo
@Desc :
*/
package entity

type TinyInfo struct {
	Id         string `form:"id" json:"id" xml:"id" bson:"_id"`
	LongUrl    string `form:"longUrl" json:"longUrl" xml:"longUrl"`
	TinyUrl    string `form:"tinyUrl" json:"tinyUrl" xml:"tinyUrl"`
	Count      int    `form:"count" json:"count" xml:"count"`
	Type       string `form:"type" json:"type" xml:"type" bson:"type"`
	CreateTime string `form:"createTime" json:"createTime" xml:"createTime"`
}
