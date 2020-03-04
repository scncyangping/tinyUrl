/*
@Time : 2019-06-14 16:41
@Author : yangping
@File : tinyInfoDao
@Desc :
*/
package mongo

import (
	"tinyUrl/common/constants"
	"tinyUrl/config"
	"tinyUrl/config/db/mongo"
	"tinyUrl/domain/entity"
)

func AddTinyInfo(tinyInfo *entity.TinyInfo) (e error) {
	var (
		err error
	)
	if err = mongo.DbInsert(config.Base.Mongo.DbName, constants.TinyInfo, tinyInfo); err != nil {
		return err
	} else {
		return nil
	}
}

func GetTinyInfoById(id string) (t entity.TinyInfo, e error) {
	var (
		tinyInfo entity.TinyInfo
		err      error
	)
	if err = mongo.DbFindById(config.Base.Mongo.DbName, constants.TinyInfo, id, &tinyInfo); err != nil {
		return tinyInfo, err
	} else {
		return tinyInfo, nil
	}
}

func AddAccessCount(id string) error {
	var (
		err error
	)

	if err = mongo.DBUpdateById(config.Base.Mongo.DbName,
		constants.TinyInfo,
		id,
		mongo.B{"$inc": mongo.B{"count": 1}}); err != nil {
		return err
	}

	return err
}
