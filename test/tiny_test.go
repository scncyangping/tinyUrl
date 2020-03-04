/*
@Time : 2019-06-15 00:50
@Author : yangping
@File : tiny_test
@Desc :
*/
package test

//
//import (
//	"fmt"
//	"go.etcd.io/etcd/pkg/testutil"
//	"testing"
//	"tinyUrl/common/util"
//	"tinyUrl/common/util/snowflake"
//)
//
///*
// * date : 2019-06-15
// * author : yangping
// * desc : 获取短链及短链反解
// */
//func TestConvert(t *testing.T) {
//	fmt.Println("testConvert start")
//	var (
//		// 默认使用62位进制转换
//		convert = util.BinaryConvert
//	)
//
//	for i := 0; i < 10000; i++ {
//		// 雪花算法获取ID
//		id := int(snowflake.NextId())
//		// 将此ID转为62进制表示
//		tinyUrl := convert.DecimalToAny(id)
//		// 将短链接转为原ID
//		nId := convert.AnyToDecimal(tinyUrl)
//
//		testutil.AssertEqual(t, id, nId)
//
//		fmt.Println("获取到的ID :", id, " 转换的短链	: ", tinyUrl, " 反转回的短链 : ", nId)
//	}
//	fmt.Println("testConvert end")
//}
//
///*
// * date : 2019-06-15
// * author : yangping
// * desc : 校验生成的短链是否有重复
// */
//func TestTinyUrl(t *testing.T) {
//	fmt.Println("testConvert start")
//	var (
//		// 默认使用62位进制转换
//		convert = util.BinaryConvert
//	)
//	mMap := make(map[string]bool, 10000)
//	num := 0
//	for i := 0; i < 10000; i++ {
//		// 雪花算法获取ID
//		id := int(snowflake.NextId())
//		// 将此ID转为62进制表示
//		tinyUrl := convert.DecimalToAny(id)
//
//		testutil.AssertFalse(t, mMap[tinyUrl])
//
//		fmt.Println(tinyUrl)
//		mMap[tinyUrl] = true
//		num++
//	}
//	fmt.Println(num)
//	fmt.Println("testConvert end")
//}
