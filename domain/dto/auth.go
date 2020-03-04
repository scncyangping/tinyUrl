/*
@date : 2020/03/03
@author : YaPi
@desc :
*/
package dto

type AuthDto struct {
	UserName string `form:"userName" json:"userName"`
	Password string `form:"password" json:"password"`
}
