/*
@Time : 2019/4/30 3:09 PM
@Author : yangping
@File : po_redis
@Desc :
*/
package confpo

type Redis struct {
	PoolSize string
	Password string
	Host     string
	Hosts    []string
}
