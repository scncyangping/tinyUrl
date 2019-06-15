#### 小结

1. 雪花算法实现分号器,将生成的int转为62进制作为短链,方便转换
2. 缓存长链 <==> 短链信息,配置过期时间,在一定时间内相同长链不生成短链,过期后同一长链可生成其他短链
3. 缓存短链 <==> 长链信息,配置过期时间,热点数据直接返回
4. 生成的短链可直接放在Redis中,设置过期时间,过期后自动失效,现直接存放于DB中,永久有效

- 时间问题,后续优化

1. 可预先生成一批短链接,有需求去取
2. 输入长、短链接校验
3. 接口为了方便调试就直接都用GET请求了,后续可按Restful风格修改

#### 使用

1. 自动生成短链接

```
GET http://localhost:9069/v1/api/tiny/url?longUrl=http://www.baidu.com


测试地址:
http://server.scncys.cn/v1/api/tiny/url?longUrl=http://www.baidu.com
```

2. 自定义短链

```
GET http://localhost:9069/v1/api/tiny/url/custom?longUrl=http://www.baidu.com&tinyUrl=test

测试地址: http://server.scncys.cn/v1/api/tiny/url/custom?longUrl=http://www.baidu.com&tinyUrl=test
```

3. 访问短链(需自己配置访问路径,这儿为了方便就这么写了)

```
http://localhost:9069/v1/api/tiny/url/go?tinyUrl=7NvKXGoh0ze

测试地址: http://server.scncys.cn/短链接
```
4. 查看短链访问次数

```
http://localhost:9069/v1/api/tiny/url/info?tinyUrl=test

测试地址: http://server.scncys.cn/v1/api/tiny/url/info?tinyUrl=test
```