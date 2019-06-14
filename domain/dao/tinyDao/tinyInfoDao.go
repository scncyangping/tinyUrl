/*
@Time : 2019-06-14 17:28
@Author : yangping
@File : tinyInfoDao
@Desc :
*/
package tinyDao

import (
	"tinyUrl/common/util"
	"tinyUrl/domain/dao/tinyDao/mongoDao"
	"tinyUrl/domain/entity"
)

func AddTinyInfo(tiny *entity.TinyInfo) (err error) {
	tiny.CreateTime = util.GetNowDateTimeFormat()
	return mongoDao.AddTinyInfo(tiny)
}

func GetTinyInfoById(id string) (t entity.TinyInfo, e error) {

	return mongoDao.GetTinyInfoById(id)
}

func AddAccessCount(id string) error {
	return mongoDao.AddAccessCount(id)
}
