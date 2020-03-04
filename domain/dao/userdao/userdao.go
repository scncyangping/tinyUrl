/*
@date : 2020/03/03
@author : YaPi
@desc :
*/
package userdao

import (
	"tinyUrl/domain/dao/userdao/mongo"
	"tinyUrl/domain/entity"
)

func AddUser(user *entity.User) error {
	return mongo.AddUser(user)
}

func GetByNameAndPassword(userName, password string) (*entity.User, error) {
	return mongo.GetByNameAndPassword(userName, password)
}

func GetByName(userName string) (*entity.User, error) {
	return mongo.GetByName(userName)
}

func GetById(id string) (*entity.User, error) {
	return mongo.GetById(id)
}
