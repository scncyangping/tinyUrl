package redis

import (
	"github.com/go-redis/redis"
	"strconv"
	"time"
	"tinyUrl/common/constants"
	"tinyUrl/common/util"
	"tinyUrl/config"
	"tinyUrl/config/log"
)

var client *redis.Client
var clusterClient *redis.ClusterClient

// client for single
func Init() error {
	var (
		host         string
		password     string
		poolSize     string
		redisOptions *redis.Options
	)

	host = config.Base.Redis.Host
	password = config.Base.Redis.Password
	poolSize = config.Base.Redis.PoolSize
	redisOptions = &redis.Options{}
	redisOptions.Addr = host
	redisOptions.DB = constants.ZERO

	if password != constants.EmptyStr {
		redisOptions.Password = password
	}
	if poolSize != constants.EmptyStr {
		redisOptions.PoolSize, _ = strconv.Atoi(poolSize)
	}

	client = redis.NewClient(redisOptions)

	_, err := client.Ping().Result()
	if err != nil {
		log.GetLogger().Errorf("Redis Init Error: Host: %v, Error:%v ", host, err)
		return err
	}

	log.GetLogger().Infof("Redis Init Success, Host: %v", host)

	return nil
}

// client for cluster
func RedisInitForCluster() error {
	var (
		hosts          []string
		password       string
		poolSize       string
		clusterOptions *redis.ClusterOptions
	)

	hosts = config.Base.Redis.Hosts
	password = config.Base.Redis.Password
	poolSize = config.Base.Redis.PoolSize

	clusterOptions = &redis.ClusterOptions{}

	clusterOptions.Addrs = hosts

	if password != constants.EmptyStr {
		clusterOptions.Password = password
	}
	if poolSize != constants.EmptyStr {
		clusterOptions.PoolSize, _ = strconv.Atoi(poolSize)
	}
	clusterClient = redis.NewClusterClient(clusterOptions)

	_, err := clusterClient.Ping().Result()
	if err != nil {
		log.GetLogger().Errorf("Redis Init Error: Hosts:%v, Error:%v ", hosts, err)
		return err
	}

	log.GetLogger().Infof("Redis Init Success, Hosts: $v", hosts)

	return nil
}

func Hset(key string, field string, value interface{}) error {
	return client.HSet(key, field, value).Err()
}

func SetTTL(key string, exTime int64) error {
	return client.Expire(key, time.Duration(exTime)*time.Duration(time.Second)).Err()
}

func Del(key string) error {
	return client.Del(key).Err()
}

func Hget(key string, field string) (string, error) {
	return client.HGet(key, field).Result()
}

func HgetAll(key string) map[string]interface{} {
	return util.ConvertStringToMap(client.HGetAll(key).Val())
}
func Get(key string) string {
	return client.Get(key).Val()
}

func SetByTtl(key string, value string, extime int64) error {
	return client.Set(key, value, time.Duration(extime)*time.Duration(time.Second)).Err()
}

func Set(key string, value string) error {
	return client.Set(key, value, time.Duration(-1)*time.Second).Err()
}

func Keys(pattern string) ([]string, error) {
	return client.Keys(pattern).Result()
}

func Scan(cursor uint64, match string, count int64) ([]string, uint64, error) {
	return client.Scan(cursor, match, count).Result()
}

func Expire(key string, extime int64) error {
	return client.Expire(key, time.Duration(extime)*time.Second).Err()
}
