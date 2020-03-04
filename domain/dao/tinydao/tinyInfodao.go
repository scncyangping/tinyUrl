/*
@Time : 2019-06-14 17:28
@Author : yangping
@File : tinyInfoDao
@Desc :
*/
package tinydao

import (
	"tinyUrl/common/util"
	"tinyUrl/domain/dao/tinydao/mongo"
	"tinyUrl/domain/entity"
)

func AddTinyInfo(tiny *entity.TinyInfo) (err error) {
	tiny.CreateTime = util.GetNowTimeStap()
	return mongo.AddTinyInfo(tiny)
}

func GetTinyInfoById(id string) (t entity.TinyInfo, e error) {

	return mongo.GetTinyInfoById(id)
}

func AddAccessCount(id string) error {
	return mongo.AddAccessCount(id)
}
