package config

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"tinyUrl/common/constants"
	"tinyUrl/config/db/confpo"
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
	Server  Server
	Log     confpo.Log
	Redis   confpo.Redis
	Mongo   confpo.Mongo
	Mysql   confpo.Mysql
	Convert confpo.BinaryConvert
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
