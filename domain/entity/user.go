/*
@date : 2020/03/03
@author : YaPi
@desc :
*/
package entity

type User struct {
	Id            string `form:"id" json:"id" xml:"id" bson:"_id"`
	UserName      string `form:"userName" json:"userName" xml:"userName"`
	Password      string `form:"password" json:"password" xml:"password"`
	Status        string `form:"status" json:"status" xml:"status"`
	Role          string `form:"role" json:"role" xml:"role"`
	CreateTime    int64  `form:"createTime" json:"createTime" xml:"createTime"`
	LastLoginTime int64  `form:"lastLoginTime" json:"lastLoginTime" xml:"lastLoginTime"`
}
