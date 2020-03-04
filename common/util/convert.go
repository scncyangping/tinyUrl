package util

import (
	"encoding/json"
	"math"
	"reflect"
	"strconv"
	"strings"
	"tinyUrl/config/log"
)

/**
 * 转换从redis获取的数据
 * @param   base {interface{}} 结构体参数
 * @returns d   {map[string]interface{}} 转换后的map
 */
func ConvertStringToMap(base map[string]string) map[string]interface{} {
	resultMap := make(map[string]interface{})
	for k, v := range base {
		var dat map[string]interface{}
		if err := json.Unmarshal([]byte(v), &dat); err == nil {
			resultMap[k] = dat
		} else {
			resultMap[k] = v
		}
	}
	return resultMap
}

/**
 * 结构体转map
 * @param   obj {interface{}} 结构体参数
 * @returns d   {map[string]interface{}} 转换后的map
 * @returns err {error} 				 错误
 */
func ConvertStructToMap(obj interface{}) (d map[string]interface{}, err error) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	d = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		d[t.Field(i).Name] = v.Field(i).Interface()
	}
	err = nil
	return
}

/**
 * date : 2019/5/7
 * author : yangping
 * desc : 结构体数据拷贝
 */
func StructCopy(DstStructPtr interface{}, SrcStructPtr interface{}) {
	srcv := reflect.ValueOf(SrcStructPtr)
	dstv := reflect.ValueOf(DstStructPtr)
	srct := reflect.TypeOf(SrcStructPtr)
	dstt := reflect.TypeOf(DstStructPtr)
	if srct.Kind() != reflect.Ptr || dstt.Kind() != reflect.Ptr ||
		srct.Elem().Kind() == reflect.Ptr || dstt.Elem().Kind() == reflect.Ptr {
		log.GetLogger().Error("Fatal error:type of parameters must be Ptr of value")
		return
	}
	if srcv.IsNil() || dstv.IsNil() {
		log.GetLogger().Error("Fatal error:value of parameters should not be nil")
		return
	}
	srcV := srcv.Elem()
	dstV := dstv.Elem()
	fields := DeepFields(reflect.ValueOf(SrcStructPtr).Elem().Type())
	for _, v := range fields {
		if v.Anonymous {
			continue
		}
		dst := dstV.FieldByName(v.Name)
		src := srcV.FieldByName(v.Name)
		if !dst.IsValid() {
			continue
		}
		if src.Type() == dst.Type() && dst.CanSet() {
			dst.Set(src)
			continue
		}
		if src.Kind() == reflect.Ptr && !src.IsNil() && src.Type().Elem() == dst.Type() {
			dst.Set(src.Elem())
			continue
		}
		if dst.Kind() == reflect.Ptr && dst.Type().Elem() == src.Type() {
			dst.Set(reflect.New(src.Type()))
			dst.Elem().Set(src)
			continue
		}
	}
	return
}

func DeepFields(baseType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField

	for i := 0; i < baseType.NumField(); i++ {
		v := baseType.Field(i)
		if v.Anonymous && v.Type.Kind() == reflect.Struct {
			fields = append(fields, DeepFields(v.Type)...)
		} else {
			fields = append(fields, v)
		}
	}

	return fields
}

var binaryConversionMap = map[int]string{
	0: "0", 1: "1", 2: "2", 3: "3", 4: "4", 5: "5",
	6: "6", 7: "7", 8: "8", 9: "9", 10: "a", 11: "b",
	12: "c", 13: "d", 14: "e", 15: "f", 16: "g", 17: "h",
	18: "i", 19: "j", 20: "k", 21: "l", 22: "m", 23: "n",
	24: "o", 25: "p", 26: "q", 27: "r", 28: "s", 29: "t",
	30: "u", 31: "v", 32: "w", 33: "x", 34: "y", 35: "z",
	36: "A", 37: "B", 38: "C", 39: "D", 40: "E", 41: "F",
	42: "G", 43: "H", 44: "I", 45: "J", 46: "K", 47: "L",
	48: "M", 49: "N", 50: "O", 51: "P", 52: "Q", 53: "R",
	54: "S", 55: "T", 56: "U", 57: "V", 58: "W", 59: "X",
	60: "Y", 61: "Z"}

var encodeURL32 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type binaryConvert struct {
	// 转换规则
	ConvertRegx string
	// 进制
	len int
}

func NewBinaryConvert(str string) *binaryConvert {
	b := &binaryConvert{}
	if len(str) < 1 {
		b.ConvertRegx = encodeURL32
	} else {
		b.ConvertRegx = str
	}
	b.len = len(b.ConvertRegx)
	return b
}

var BinaryConvert = NewBinaryConvert(encodeURL32)

func (b *binaryConvert) DecimalToAny(num int) string {
	n := b.len
	newNumStr := ""
	var remainder int
	var remainderString string
	for num != 0 {
		remainder = num % n
		if b.len+1 > remainder && remainder > 9 {
			// remainderString = binaryConversionMap[remainder]
			remainderString = string(encodeURL32[remainder])
		} else {
			remainderString = strconv.Itoa(remainder)
		}
		newNumStr = remainderString + newNumStr
		num = num / n
	}
	return newNumStr
}

func binaryConversionKey(in string) int {
	result := -1
	for k, v := range binaryConversionMap {
		if in == v {
			result = k
		}
	}
	return result
}

func (b *binaryConvert) AnyToDecimal(num string) int {
	var newNum int
	nNum := len(strings.Split(num, "")) - 1
	for _, value := range strings.Split(num, "") {
		tmp := binaryConversionKey(value)
		if tmp != -1 {
			newNum = newNum + tmp*int(math.Pow(float64(b.len), float64(nNum)))
			nNum--
		} else {
			break
		}
	}
	return newNum
}
