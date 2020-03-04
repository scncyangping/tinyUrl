package util

import (
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"reflect"
	"tinyUrl/domain/dto"
)

const (
	Nil = iota
	Bool
	Int
	Float
	Array
	Chan
	Func
	Interface
	Map
	Ptr
	Slice
	String
	Struct
	Undefined
)

func IsExpectType(i interface{}) int {
	if i == nil {
		return Nil
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.String:
		return String
	case reflect.Map:
		return Map
	case reflect.Array:
		return Array
	case reflect.Struct:
		return Struct
	case reflect.Slice:
		return Slice
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Int
	case reflect.Float32, reflect.Float64:
		return Float
	case reflect.Interface:
		return Interface
	case reflect.Bool:
		return Bool
	default:
		return Undefined
	}
}

/**
 * 判断obj是否在target中，target支持的类型arrary,slice,map
 */
func Contains(obj interface{}, target interface{}) bool {
	if target == nil {
		return false
	}
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

/**
 * 深度拷贝
 */
func DeepCopy(value interface{}) interface{} {
	if valueMap, ok := value.(map[string]interface{}); ok {
		newMap := make(map[string]interface{})
		for k, v := range valueMap {
			newMap[k] = DeepCopy(v)
		}

		return newMap
	} else if valueSlice, ok := value.([]interface{}); ok {
		newSlice := make([]interface{}, len(valueSlice))
		for k, v := range valueSlice {
			newSlice[k] = DeepCopy(v)
		}

		return newSlice
	} else if valueMap, ok := value.(bson.M); ok {
		newMap := make(bson.M)
		for k, v := range valueMap {
			newMap[k] = DeepCopy(v)
		}
	}
	return value
}

/**
 * 组装需要使用的参数,过滤不需要的字段
 * @param   data 	原数据
 * @param field   使用字段
 */
func FiltrationData(data map[string]interface{}, field []string) map[string]interface{} {
	resultMap := make(map[string]interface{})
	if data == nil || field == nil || len(data) < 1 || len(field) < 1 {
		return nil
	}

	for i := 0; i < len(field); i++ {
		for j := 0; j < len(data); j++ {
			if Contains(field[i], data) {
				resultMap[field[i]] = data[field[i]]
			}
		}
	}
	return resultMap
}

func GetSession(ctx *gin.Context) *dto.Session {
	var (
		session *dto.Session
	)
	if s, ok := ctx.Get("session"); ok {
		if session, ok = s.(*dto.Session); ok {
			return session
		}
	}
	return nil
}
