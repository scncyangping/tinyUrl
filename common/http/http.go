package http

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"reflect"
	"strconv"
	"tinyUrl/common/constants"
	"tinyUrl/common/util"
	"tinyUrl/config/log"
)

const (
	DefaultStatus             = -1
	StatusOK                  = 200
	StatusMovedPermanently    = 301
	StatusFound               = 302
	StatusBadRequest          = 400
	StatusNotFound            = 404
	StatusInternalServerError = 500

	// Request Error
	RequestParameterError    = 1001
	RequestCheckTokenError   = 1002
	RequestCheckTokenTimeOut = 1003
	RequestTokenNotFount     = 1004
	CreateTokenError         = 1005
	// System  Error
	DataConvertError      = 2001
	ParameterConvertError = 2002
	// DataBase Error
	InitDataBaseError = 3001
	QueryDBError      = 3002
	UserNotFound      = 3003
	AddUserError      = 3004
)

const (
	MethodGet     = "GET"
	MethodHead    = "HEAD"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH"
	MethodDelete  = "DELETE"
	MethodConnect = "CONNECT"
	MethodOptions = "OPTIONS"
	MethodTrace   = "TRACE"

	ConnectTypeJson = "application/json"
	ConnectTypeWww  = "application/x-www-form-urlencoded"
)

var statusText = map[int]string{
	DefaultStatus:             "",
	StatusOK:                  "OK",
	StatusBadRequest:          "Bad Request",
	StatusMovedPermanently:    "Moved Permanently",
	StatusFound:               "Found",
	StatusNotFound:            "Not Found",
	StatusInternalServerError: "Internal Server Error",
	RequestParameterError:     "Request Parameter Error",
	DataConvertError:          "Data Convert Error",
	RequestCheckTokenError:    "Token Is Not Exists, Please Login",
	ParameterConvertError:     "Parameter Error, Please Check Parameter",
	InitDataBaseError:         "Init DataBase Error",
	QueryDBError:              "Query DataBase Error",
	RequestCheckTokenTimeOut:  "request check token time out",
	RequestTokenNotFount:      "request token not found, please login first",
	UserNotFound:              "user not found",
	CreateTokenError:          "create token error",
	AddUserError:              "add user error",
}

// StatusText returns a text for the HTTP status code. It returns the empty
// string if the code is unknown.
func StatusText(code int) string {
	return statusText[code]
}

type QuerySet struct {
	Query map[string]interface{}
	Skip  int
	Limit int
	Sort  string
}

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (res *Result) Set(code int, msg string) {
	res.Code = code
	res.Msg = msg
}

func (res *Result) SetMsg(msg string) {
	res.Msg = msg
}

func (res *Result) SetCode(code int) {
	res.Code = code
}

func Instance() *Result {
	return &Result{DefaultStatus, StatusText(DefaultStatus), nil}
}

func SendSuccess(ctx *gin.Context) {
	rep := Instance()
	rep.SetCode(StatusOK)
	rep.SetMsg(StatusText(StatusOK))
	SendSuccessRep(ctx, rep)
}

func SendSuccessRep(ctx *gin.Context, rep *Result) {
	var (
		// 当前请求ID
		actionId string
	)
	if rep.Msg == constants.EmptyStr {
		rep.Msg = StatusText(StatusOK)
	}

	rep.Set(StatusOK, rep.Msg)

	actionId = ctx.GetString("ActionId")

	log.GetLogger().Infof("Response [%s] To Client: %v", actionId, rep)
	ctx.JSON(StatusOK, rep)
}

func SendFailure(ctx *gin.Context) {
	rep := Instance()
	rep.SetCode(StatusBadRequest)
	rep.SetMsg(StatusText(StatusBadRequest))
	SendFailureRep(ctx, rep)
}

func SendFailureRep(ctx *gin.Context, rep *Result) {
	var (
		// 当前请求ID
		actionId string
	)

	if rep.Code == DefaultStatus {
		rep.Code = StatusBadRequest
	}

	rep.Msg = StatusText(rep.Code)

	actionId = ctx.GetString("ActionId")

	log.GetLogger().Infof("Response [%s] To Client: %v", actionId, rep)

	ctx.JSON(StatusBadRequest, rep)
}

func SendFailureError(ctx *gin.Context, rep *Result, err error) {
	var (
		// 当前请求ID
		actionId string
	)

	if rep.Code == constants.ZERO {
		rep.Code = StatusBadRequest
	}

	if rep.Msg == constants.EmptyStr {
		rep.Msg = StatusText(rep.Code)
	}

	actionId = ctx.GetString("ActionId")

	if err != nil {
		log.GetLogger().Errorf("Response [%s] to Client: %v ,Error is %v", actionId, rep, err)
	} else {
		log.GetLogger().Infof("Response [%s] To Client: %v", actionId, rep)
	}
	ctx.JSON(StatusBadRequest, rep)
}

/**
构造新增、更新用的db对象，适用于不适合用struct组装请求数据的情况
@param
obj: 结构体对象,如果数据中有数字类型的字段，会按结构体指定的类型进行转换 ；不需要的时候传nil
@return
updateSet : 更新db用的map
ok： true，处理数据成功； false ，处理失败
*/
func UpdateData(ctx *gin.Context, obj interface{}) (updateSet map[string]interface{}, ok bool) {
	var (
		err     error
		context *Result
		params  map[string]interface{}
	)

	context = Instance()
	params = map[string]interface{}{}

	method := ctx.Request.Method

	//contentType := ctx.Request.Header.Get("Content-Type")
	if method == MethodGet {
		return
	} else {
		err = ctx.Bind(&params)
		if err != nil {
			context.SetCode(RequestParameterError)
			context.SetMsg(StatusText(RequestParameterError))
			SendFailureError(ctx, context, err)
			return
		}
	}

	if obj != nil {
		for k, v := range params {
			vStr := convertDouble2String(v)
			if vStr != "" {
				err = convertValType(k, vStr, obj, params)
				if err != nil {
					context.SetCode(DataConvertError)
					context.SetMsg(StatusText(DataConvertError))
					SendFailureError(ctx, context, err)
					return
				}
			}
		}
	}

	updateSet = params
	ok = true

	log.GetLogger().Infof("Get Request Info: [%v]", updateSet)
	return
}

/**
构造查询用db对象，支持json、form、query格式的参数
@param
obj: 结构体对象,如果传值则按结构体指定的字段类型进行转换 ；不需要的时候传nil
@return
queryObj : 查询db用的对象，增加了分页用的值
ok： true，处理数据成功； false ，处理失败
*/
func QueryDataForStuct(ctx *gin.Context, obj interface{}) (queryObj QuerySet, ok bool) {

	var (
		context     *Result
		contentType string
		query       *map[string]interface{}
		page        = constants.DefaultPage
		limit       = constants.DefaultPageSize
	)

	context = Instance()
	contentType = ctx.Request.Header.Get("Content-Type")
	query = &map[string]interface{}{}

	if contentType == ConnectTypeJson {
		params := make(map[string]interface{})
		err := ctx.BindJSON(&params)
		if err != nil {
			context.SetCode(RequestParameterError)
			context.SetMsg(StatusText(RequestParameterError))
			SendFailureError(ctx, context, err)
			return
		}
		for k, v := range params {
			if v == nil || v == "" || v == 0 {
				delete(params, k)
			} else {
				if k == "page" {
					if util.IsExpectType(v) == util.String {
						page, _ = strconv.Atoi(v.(string))
					} else {
						page = int(v.(float64))
					}
				} else if k == "pageSize" {
					if util.IsExpectType(v) == util.String {
						limit, _ = strconv.Atoi(v.(string))
					} else {
						limit = int(v.(float64))
					}
				} else {
					err := convertValTypeData(k, v, obj, *query)
					if err != nil {
						context.SetCode(DataConvertError)
						context.SetMsg(StatusText(DataConvertError))
						SendFailureError(ctx, context, err)
						return
					}
				}
			}
		}
	} else {
		err := ctx.Request.ParseForm()
		if err != nil {
			context.SetCode(DataConvertError)
			context.SetMsg(StatusText(DataConvertError))
			SendFailureError(ctx, context, err)
			return
		}
		params := ctx.Request.Form
		for k, v := range params {
			l := len(v)
			if l == 1 && v[0] == "" {
				delete(params, k)
			} else {
				if k == "page" {
					page, _ = strconv.Atoi(v[0])
				} else if k == "pageSize" {
					limit, _ = strconv.Atoi(v[0])
				} else if k == "sessionId" {
					continue
				} else {
					if l > 1 {
					} else {
						err := convertValType(k, v[0], obj, *query)
						if err != nil {
							context.SetCode(DataConvertError)
							context.SetMsg(StatusText(DataConvertError))
							SendFailureError(ctx, context, err)
							return
						}
					}
				}
			}
		}
	}

	skip := 0
	if page > 0 {
		skip = (page - 1) * limit
	}
	queryObj.Query = *query
	queryObj.Skip = skip
	queryObj.Limit = limit
	ok = true

	return
}

/**
根据struct反射转换map的数据类型
*/
func convertValType(k string, v string, obj interface{}, maps map[string]interface{}) error {
	typ := reflect.TypeOf(obj).Elem()
	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		fieldName := typeField.Tag.Get("json")
		if fieldName == k {
			if err := setWithProperType(k, typeField.Type.Kind(), v, maps); err != nil {
				return err
			}
			break
		}
	}
	return nil
}
func setWithProperType(key string, valueKind reflect.Kind, val string, maps map[string]interface{}) error {
	switch valueKind {
	case reflect.Int:
		return setIntField(key, val, 0, maps)
	case reflect.Int8:
		return setIntField(key, val, 8, maps)
	case reflect.Int16:
		return setIntField(key, val, 16, maps)
	case reflect.Int32:
		return setIntField(key, val, 32, maps)
	case reflect.Int64:
		return setIntField(key, val, 64, maps)
	case reflect.Uint:
		return setUintField(key, val, 0, maps)
	case reflect.Uint8:
		return setUintField(key, val, 8, maps)
	case reflect.Uint16:
		return setUintField(key, val, 16, maps)
	case reflect.Uint32:
		return setUintField(key, val, 32, maps)
	case reflect.Uint64:
		return setUintField(key, val, 64, maps)
	case reflect.Float32:
		return setFloatField(key, val, 32, maps)
	case reflect.Float64:
		return setFloatField(key, val, 64, maps)
	case reflect.Bool:
		return setBoolField(key, val, maps)
	case reflect.String:
		maps[key] = val
	default:
		return errors.New("Unknown type")
	}
	return nil
}

func setIntField(key string, val string, bitSize int, maps map[string]interface{}) error {
	if val == "" {
		val = "0"
	}
	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil {
		switch bitSize {
		case 0:
			maps[key] = int(intVal)
		case 8:
			maps[key] = int8(intVal)
		case 16:
			maps[key] = int16(intVal)
		case 32:
			maps[key] = int32(intVal)
		case 64:
			maps[key] = int64(intVal)
		default:
			return errors.New("Unknown type")
		}
	}
	return err
}

func setUintField(key string, val string, bitSize int, maps map[string]interface{}) error {
	if val == "" {
		val = "0"
	}
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		switch bitSize {
		case 0:
			maps[key] = uint(uintVal)
		case 8:
			maps[key] = uint8(uintVal)
		case 16:
			maps[key] = uint16(uintVal)
		case 32:
			maps[key] = uint32(uintVal)
		case 64:
			maps[key] = uint64(uintVal)
		default:
			return errors.New("Unknown type")
		}
	}
	return err
}

func setBoolField(key string, val string, maps map[string]interface{}) error {
	if val == "" {
		val = "false"
	}
	boolVal, err := strconv.ParseBool(val)
	if err == nil {
		maps[key] = boolVal
	}
	return err
}

func setFloatField(key string, val string, bitSize int, maps map[string]interface{}) error {
	if val == "" {
		val = "0.0"
	}
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err == nil {
		switch bitSize {
		case 32:
			maps[key] = float32(floatVal)
		case 64:
			maps[key] = float64(floatVal)
		default:
			return errors.New("Unknown type")
		}
	}
	return err
}

func convertValTypeData(k string, v interface{}, obj interface{}, maps map[string]interface{}) error {
	typ := reflect.TypeOf(obj).Elem()
	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		fieldName := typeField.Tag.Get("form")
		if fieldName == k {
			if err := setWithProperTypeData(k, typeField.Type.Kind(), v, maps); err != nil {
				return err
			}
			break
		}
	}
	return nil
}

/**
js传输过来的数字类型默认是double，需要按需进行转换
如果是数字类型，统一转换为string返回，返回为""说明是非数字类型
*/
func convertDouble2String(v interface{}) string {
	v1, ok1 := v.(float64)
	var vStr string
	if ok1 {
		vStr = strconv.FormatFloat(v1, 'f', -1, 64)
	}

	return vStr
}

func setWithProperTypeData(key string, valueKind reflect.Kind, val interface{}, maps map[string]interface{}) error {
	vStr := convertDouble2String(val)
	if vStr != "" {
		return setWithProperType(key, valueKind, vStr, maps)
	} else {
		switch valueKind {
		case reflect.Bool:
			maps[key] = val.(bool)
		case reflect.String:
			maps[key] = val
		case reflect.Map:
			maps[key] = val
		case reflect.Array:
			maps[key] = val
		case reflect.Slice:
			maps[key] = val
		case reflect.Struct:
			maps[key] = val
		default:
			return errors.New("Unknown type")
		}
		return nil
	}
}

func GetLimitAndSkip(data map[string]interface{}) (int, int) {
	page := 1
	limit := 10
	skip := 0
	if util.Contains("page", data) {
		if util.IsExpectType(data["page"]) == util.String {
			page, _ = strconv.Atoi(data["page"].(string))
		} else {
			_, ok := data["page"].(float64)

			if !ok {
				page = data["page"].(int)
			} else {
				page = int(data["page"].(float64))
			}
		}
		delete(data, "page")
	}

	if util.Contains("pageSize", data) {
		if util.IsExpectType(data["pageSize"]) == util.String {
			limit, _ = strconv.Atoi(data["pageSize"].(string))
		} else {
			_, ok := data["pageSize"].(float64)
			if !ok {
				limit, _ = data["pageSize"].(int)
			} else {
				limit = int(data["pageSize"].(float64))
			}

		}
		delete(data, "pageSize")
	}
	if page > 0 {
		skip = (page - 1) * limit
	}
	return skip, limit
}

type RData struct {
	Data  map[string]interface{}
	Limit int
	Skip  int
}

// 获取请求参数
func QueryData(ctx *gin.Context) RData {
	queryData, err := ctx.GetRawData()
	if err != nil {
		log.GetLogger().Info(err)
	}
	dat := make(map[string]interface{})
	err = json.Unmarshal(queryData, &dat)
	if err != nil {
		log.GetLogger().Info(err)
	}
	resData, _ := ClearNilField(dat)
	result, ok := resData.(map[string]interface{})
	if !ok {
		log.GetLogger().Infof(StatusText(DataConvertError))
	}
	skip, limit := GetLimitAndSkip(result)

	return RData{
		Data:  result,
		Skip:  skip,
		Limit: limit,
	}
}

// 为空的字段默认不拷贝,不传递
// 若为数组或slice类型,长度为0,不拷贝
// float64类型转int
func ClearNilField(value interface{}) (interface{}, bool) {
	if valueMap, ok := value.(map[string]interface{}); ok {
		newMap := make(map[string]interface{})
		for k, v := range valueMap {
			if value, res := ClearNilField(v); res {
				newMap[k] = value
			}
		}
		if len(newMap) < constants.ONE {
			return nil, false
		} else {
			return newMap, true
		}
	} else if valueSlice, ok := value.([]interface{}); ok {
		newSlice := make([]interface{}, constants.ZERO)
		for _, v := range valueSlice {
			if value, res := ClearNilField(v); res {
				newSlice = append(newSlice, value)
			}
		}
		if len(newSlice) < constants.ONE {
			return nil, false
		} else {
			return newSlice, true
		}
	}
	if util.IsExpectType(value) == util.Float {
		v1, ok1 := value.(float64)
		if ok1 {
			return int(v1), true
		}

	}
	if util.IsExpectType(value) == util.String {
		if value == constants.EmptyStr {
			return nil, false
		}
	}
	return value, true
}
