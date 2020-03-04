/*
@date : 2020/03/03
@author : YaPi
@desc :
*/
package entity

type JwtToken struct {
	Id         string `form:"id" json:"id" xml:"id" bson:"_id"`
	UserName   string `form:"userName" json:"userName" xml:"userName"`
	Password   string `form:"password" json:"password" xml:"password"`
	CreateTime int64  `form:"createTime" json:"createTime" xml:"createTime"`
}
