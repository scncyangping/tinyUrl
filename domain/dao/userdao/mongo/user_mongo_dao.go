/*
@date : 2020/03/03
@author : YaPi
@desc :
*/
package mongo

import (
	"tinyUrl/common/constants"
	"tinyUrl/config"
	"tinyUrl/config/db/mongo"
	"tinyUrl/domain/entity"
)

func AddUser(user *entity.User) error {
	var (
		err error
	)
	if err = mongo.DbInsert(config.Base.Mongo.DbName, constants.User, user); err != nil {
		return err
	} else {
		return nil
	}
}

func GetByName(userName string) (*entity.User, error) {
	var (
		u   entity.User
		err error
	)
	query := mongo.B{"name": userName}
	if err = mongo.DBFind(config.Base.Mongo.DbName, constants.User, query, &u); err != nil {
		return &u, err
	} else {
		return &u, nil
	}
}

func GetByNameAndPassword(userName, password string) (*entity.User, error) {
	var (
		u   *entity.User
		err error
	)
	query := mongo.B{"username": userName, "password": password}
	if err = mongo.DbFindOne(config.Base.Mongo.DbName, constants.User, query, &u); err != nil {
		return u, err
	} else {
		return u, nil
	}
}

func GetById(id string) (*entity.User, error) {
	var (
		u   entity.User
		err error
	)
	if err = mongo.DbFindById(config.Base.Mongo.DbName, constants.User, id, &u); err != nil {
		return &u, err
	} else {
		return &u, nil
	}
}
