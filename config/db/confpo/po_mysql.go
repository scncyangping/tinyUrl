/*
@Time : 2019/4/30 3:09 PM
@Author : yangping
@File : po_mysql
@Desc :
*/
package confpo

type Mysql struct {
	// 连接用户名
	User string
	// 连接密码
	Password string
	// 连接地址
	Host string
	// 连接地址数组 集群连接时使用
	Hosts []string
	// 数据库名称
	DbName string `yaml:"dbName"`
	// 连接参数,配置编码、是否启用ssl等 charset=utf8
	ConnectInfo string
	// 用于设置最大打开的连接数，默认值为0表示不限制
	MaxOpenConn int
	// 用于设置闲置的连接数
	MaxIdleConn int
	// 连接名称 用于多个连接时区分
	ConnName string
}
