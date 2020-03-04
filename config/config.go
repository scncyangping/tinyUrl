package config

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"tinyUrl/common/constants"
	"tinyUrl/config/db/confpo"
	"tinyUrl/config/global/global_po"
)

var (
	configYml     = "./config.yml"
	confMap       = make(map[string]string)
	consulConfMap = make(map[string]string)
	Base          = baseConfig{}
)

// 系统配置
type Server struct {
	Name string
	Port string
}

type baseConfig struct {
	// 服务设置
	Server Server
	// 日志设置
	Log confpo.Log
	// redis设置
	Redis confpo.Redis
	// mongodb 设置
	Mongo confpo.Mongo
	// mysql 设置
	Mysql confpo.Mysql
	// 转换参数设置
	Convert confpo.BinaryConvert
	// Jwt设置
	Jwt global_po.Jwt
}

func SetConfFile(file string) {
	configYml = file
}

func GetConf(key string) string {
	return confMap[key]
}

func Init() {
	f := flag.String("f", constants.EmptyStr, constants.StartMessage)
	flag.Parse()
	if *f != constants.EmptyStr {
		SetConfFile(*f)
	} else {
		// 判断环境变量
		mode := os.Getenv(constants.StartRunMode)

		if mode == constants.PRO {
			SetConfFile(constants.ProConfigFile)
		} else {
			SetConfFile(constants.DevConfigFile)
		}
	}

	conf, err := ioutil.ReadFile(configYml)
	if err != nil {
		panic(err)
		return
	}

	InitYml(conf)

}

// 加载yml配置文件
func InitYml(byteArray []byte) {
	err := yaml.Unmarshal(byteArray, &Base)
	if err != nil {
		panic(err)
	}
}

// 加载普通配置文件
func InitConf(byteArray []byte) {
	confStr := string(byteArray)
	confStrSlice := strings.Split(confStr, "\n")
	for i := 0; i < len(confStrSlice); i++ {
		confStrSlice[i] = strings.Trim(confStrSlice[i], "\r")
		if strings.HasPrefix(confStrSlice[i], "//") || confStrSlice[i] == "" {
			continue
		}
		oneConf := strings.Split(confStrSlice[i], "=")
		if len(oneConf) == constants.TWO {
			confMap[strings.Trim(oneConf[constants.ZERO], "\r")] =
				strings.Trim(oneConf[constants.ONE], "\r")
		}
	}
}
