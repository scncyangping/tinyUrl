/*
@date : 2020/03/03
@author : YaPi
@desc :
*/
package dto

type Session struct {
	UserName  string `form:"userName" json:"userName" xml:"userName"`
	Password  string `form:"password" json:"password" xml:"password"`
	Ip        string `form:"ip" json:"ip" xml:"ip"`
	LoginType string `form:"loginType" json:"loginType" xml:"loginType"`
}
