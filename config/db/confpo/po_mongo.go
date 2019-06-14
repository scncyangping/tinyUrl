/*
@Time : 2019/4/30 3:05 PM
@Author : yangping
@File : po
@Desc :
*/
package confpo

type Mongo struct {
	Host     string
	User     string
	DbName   string `yaml:"dbName"`
	Password string
	PoolSize int
}
