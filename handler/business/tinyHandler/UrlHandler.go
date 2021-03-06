/*
@Time : 2019-06-14 10:17
@Author : yangping
@File : UrlHandler
@Desc :
*/
package tinyHandler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"tinyUrl/common/constants"
	"tinyUrl/common/http"
	"tinyUrl/common/util"
	"tinyUrl/common/util/snowflake"
	"tinyUrl/config"
	"tinyUrl/config/db/redis"
	"tinyUrl/domain/dao/tinyDao"
	"tinyUrl/domain/dto"
	"tinyUrl/domain/entity"
	"tinyUrl/domain/vo"
)

/*
 * date : 2019-06-15
 * author : yangping
 * desc : 对应短链获取计数
 */
func UrlBaseInfo(ctx *gin.Context) {
	var (
		tinyDto dto.TinyDto
		convert = util.NewBinaryConvert(config.Base.Convert.BinaryStr)
		err     error
		// 初始化返回结构体
		result = http.Instance()
	)
	// 请求参数校验
	if err = ctx.Bind(&tinyDto); err != nil ||
		tinyDto.TinyUrl == constants.EmptyStr {
		result.Code = http.ParameterConvertError
		http.SendFailureRep(ctx, result)
		return
	}

	// 若 Redis 不存在此key, 查询DB内是否有对应key
	tinyId := strconv.Itoa(convert.AnyToDecimal(tinyDto.TinyUrl))

	if t, err := tinyDao.GetTinyInfoById(tinyId); err != nil {
		result.Code = http.QueryDBError
		http.SendFailureRep(ctx, result)
	} else {
		result.Data = &vo.TinyVO{
			LongUrl: t.LongUrl,
			TinyUrl: t.TinyUrl,
			Count:   t.Count,
		}
		http.SendSuccessRep(ctx, result)
	}
}

/*
 * date : 2019-06-14
 * author : yangping
 * desc : 转换长网址到短网址
 */
func UrlTransform(ctx *gin.Context) {
	var (
		tinyInfo entity.TinyInfo
		tinyDto  dto.TinyDto
		err      error
		// 短连接
		tinyUrl string
		// 雪花算法生成ID
		id = int(snowflake.NextId())
		// 获取进制转换工具
		convert = util.NewBinaryConvert(config.Base.Convert.BinaryStr)
		// 初始化返回结构体
		result = http.Instance()
	)

	// 请求参数校验
	if err = ctx.Bind(&tinyDto); err != nil || tinyDto.LongUrl == constants.EmptyStr {
		result.Code = http.ParameterConvertError
		http.SendFailureRep(ctx, result)
		return
	}

	// 查询此短链是否存在 存在直接返回 -- Redis
	if isExist, _ := checkLongUrl(tinyDto.LongUrl); isExist {
		result.Data = "此长链已存在"
		http.SendSuccessRep(ctx, result)
		return
	}
	// 相同长链对应多个短链, 若需要 1对1, 单独校重处理
	// 将ID转化为62进制
	tinyUrl = convert.DecimalToAny(id)

	tinyInfo.LongUrl = tinyDto.LongUrl
	tinyInfo.Count = constants.ZERO
	tinyInfo.Id = strconv.Itoa(id)
	tinyInfo.TinyUrl = tinyUrl
	tinyInfo.Type = constants.ConvertDefault

	if err = tinyDao.AddTinyInfo(&tinyInfo); err != nil {
		http.SendFailureError(ctx, result, err)
	} else {
		// 放在Redis中
		// 长链Redis中
		addLongUrlRedisKey(tinyInfo.LongUrl, tinyInfo.TinyUrl, tinyInfo.Id)

		// 短链放Redis中
		addTinyUrlRedisKey(tinyInfo.TinyUrl, tinyInfo.LongUrl, tinyInfo.Id)

		result.Data = &dto.TinyDto{
			LongUrl: tinyInfo.LongUrl,
			TinyUrl: tinyInfo.TinyUrl,
		}
		http.SendSuccessRep(ctx, result)
	}
}

/*
 * date : 2019-06-14
 * author : yangping
 * desc : 自定义短链接key
 */
func UrlTransformCustom(ctx *gin.Context) {
	var (
		tinyInfo entity.TinyInfo
		tinyDto  dto.TinyDto
		err      error
		// 初始化返回结构体
		result  = http.Instance()
		convert = util.NewBinaryConvert(config.Base.Convert.BinaryStr)
	)

	// 请求参数校验
	if err = ctx.Bind(&tinyDto); err != nil ||
		tinyDto.LongUrl == constants.EmptyStr ||
		tinyDto.TinyUrl == constants.EmptyStr {
		result.Code = http.ParameterConvertError
		http.SendFailureRep(ctx, result)
		return
	}

	// 查询此长链是否存在 存在直接返回 -- Redis
	if isExist, _ := checkLongUrl(tinyDto.LongUrl); isExist {
		result.Data = "此长链已存在"
		http.SendSuccessRep(ctx, result)
		return
	}

	// 查询此短链是否存在 存在直接返回 -- Redis
	// 自定义短链接会校验DB
	if isExist, _ := checkTinyUrl(tinyDto.TinyUrl, true); isExist {
		result.Data = "此短链已存在"
		http.SendSuccessRep(ctx, result)
		return
	}

	// 不存在就新增
	tinyInfo.Id = strconv.Itoa(convert.AnyToDecimal(tinyDto.TinyUrl))
	tinyInfo.LongUrl = tinyDto.LongUrl
	tinyInfo.Count = constants.ZERO
	tinyInfo.TinyUrl = tinyDto.TinyUrl
	tinyInfo.Type = constants.ConvertCustom

	if err = tinyDao.AddTinyInfo(&tinyInfo); err != nil {
		http.SendFailureError(ctx, result, err)
	} else {
		// 长链Redis中
		addLongUrlRedisKey(tinyInfo.LongUrl, tinyInfo.TinyUrl, tinyInfo.Id)

		// 短链放Redis中
		addTinyUrlRedisKey(tinyInfo.TinyUrl, tinyInfo.LongUrl, tinyInfo.Id)

		result.Data = &dto.TinyDto{
			LongUrl: tinyInfo.LongUrl,
			TinyUrl: tinyInfo.TinyUrl,
		}

		http.SendSuccessRep(ctx, result)
	}
}

/*
 * date : 2019-06-14
 * author : yangping
 * desc : 短链接跳转长链接
 */
func Redirect4TinyUrl(ctx *gin.Context) {
	var (
		tinyDto dto.TinyDto
		err     error
		// 初始化返回结构体
		result = http.Instance()
	)
	// 请求参数校验
	if err = ctx.Bind(&tinyDto); err != nil ||
		tinyDto.TinyUrl == constants.EmptyStr {
		result.Code = http.ParameterConvertError
		http.SendFailureRep(ctx, result)
		return
	}
	// 查询此短链是否存在 存在直接返回 -- Redis
	// 自定义短链接会校验DB
	if isExist, longUrl := checkTinyUrl(tinyDto.TinyUrl, true); isExist {

		// 若对应长链接存在,需要统计访问信息
		// 可以放在消息队列里面去做 便于更多样的统计
		// 这儿直接单开线程 同步信息到DB中
		array := strings.Split(longUrl, constants.UnderLine)

		go tinyDao.AddAccessCount(array[len(array)-1])

		ctx.Redirect(http.StatusFound, array[constants.ZERO])
		return
	}
	result.Data = "url not found"
	http.SendFailureRep(ctx, result)
}

/*
 * date : 2019-06-14
 * author : yangping
 * desc : 校验Redis是否存在此长链接对应key 可设置是否校验DB
 */
func checkLongUrl(longUrl string) (bool, string) {
	var (
		redisKey string
		str      string
		err      error
	)

	redisKey = fmt.Sprintf("%s:%s:%s", constants.SCNCYS, constants.LongUrl, longUrl)

	// 查询此短链是否存在 存在直接返回 -- Redis
	if str, err = getRedisKey(redisKey, false); err == nil {
		return true, str
	}
	return false, constants.EmptyStr
}

/*
 * date : 2019-06-14
 * author : yangping
 * desc : 同一长链接在设置时间周期内,不允许重复生成短链接,防止攻击
 */
func addLongUrlRedisKey(longUrl, tinyUrl, id string) (bool, error) {
	var (
		redisKey string
		str      string
		err      error
	)

	redisKey = fmt.Sprintf("%s:%s:%s", constants.SCNCYS, constants.LongUrl, longUrl)
	// 将这一条记录放在Redis当中

	str = tinyUrl + constants.UnderLine + id
	err = redis.SetByTtl(redisKey, str, config.Base.Convert.LongUrlExpire)

	if err != nil {
		return false, err
	}
	return true, nil
}

/*
 * date : 2019-06-14
 * author : yangping
 * desc : 添加短链缓存
 */
func addTinyUrlRedisKey(tinyUrl, longUrl, id string) (bool, error) {
	var (
		redisKey string
		str      string
		err      error
	)

	redisKey = fmt.Sprintf("%s:%s:%s", constants.SCNCYS, constants.TinyUrl, tinyUrl)
	// 将这一条记录放在Redis当中

	str = longUrl + constants.UnderLine + id
	err = redis.SetByTtl(redisKey, str, config.Base.Convert.TinyUrlExpire)

	if err != nil {
		return false, err
	}
	return true, nil
}

/*
 * date : 2019-06-14
 * author : yangping
 * desc : 校验Redis是否存在此短链接对应key 可设置是否校验DB
 */
func checkTinyUrl(tinyUrl string, checkDb bool) (bool, string) {
	var (
		convert  = util.NewBinaryConvert(config.Base.Convert.BinaryStr)
		redisKey string
	)

	redisKey = fmt.Sprintf("%s:%s:%s", constants.SCNCYS, constants.TinyUrl, tinyUrl)

	// 查询此短链是否存在 存在直接返回 -- Redis
	if str, err := getRedisKey(redisKey, true); err == nil {
		return true, str
	} else {
		if checkDb {
			// 若 Redis 不存在此key, 查询DB内是否有对应key
			tinyId := strconv.Itoa(convert.AnyToDecimal(tinyUrl))
			t, error := tinyDao.GetTinyInfoById(tinyId)
			// 将这一条记录放在Redis当中
			if error == nil {
				addTinyUrlRedisKey(t.TinyUrl, t.LongUrl, t.Id)
				str := t.LongUrl + constants.UnderLine + t.Id
				return true, str
			}
		}
	}
	return false, constants.EmptyStr
}

/*
 * date : 2019-06-14
 * author : yangping
 * desc : 获取Redis的值,若存在则更新过期时间
 */
func getRedisKey(redisKey string, upExpire bool) (string, error) {
	var (
		tinyUrl string
	)

	tinyUrl = redis.Get(redisKey)

	if tinyUrl != constants.EmptyStr {

		if upExpire {
			if err := redis.Expire(redisKey, constants.ExpireTime); err != nil {
				return constants.EmptyStr, errors.New("update expire time error")
			}
		}
		return tinyUrl, nil

	} else {
		return constants.EmptyStr, errors.New("value does not exist")
	}
}
