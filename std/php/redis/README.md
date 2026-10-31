# Redis 可选模块

提供全局 `Redis`、`RedisException` 及常用 phpredis 接口，使用 Go 网络库直接实现
RESP2 协议。模块版本为 `0.1.0-origami`，不冒充完整 phpredis 6.x 扩展。

## 显式注册

`std.Load`、`php.Load`、`framework.Load` 和默认 CLI 不加载此驱动。
在需要 Redis 的宿主 `main.go` 的实际 VM 初始化代码中加入：

```go
import (
    "github.com/php-any/origami/std/php"
    "github.com/php-any/origami/std/php/redis"
)

// 创建 VM 并加载基础标准库后：
php.Load(vm)
redis.Load(vm)
```

未注册时 `extension_loaded('redis')` 返回 false。声明同名的用户类或 Laravel
Redis facade 也不会让扩展显示为已加载。注册过程不连接服务器。

```php
<?php
$redis = new Redis();
$redis->connect('127.0.0.1', 6379, 2.0, null, 0, 2.0);
$redis->setOption(Redis::OPT_PREFIX, 'myapp:');
$redis->set('hello', 'world', ['ex' => 60]);
echo $redis->get('hello');

$results = $redis->pipeline()->set('counter', 0)->incr('counter')->get('counter')->exec();
// [true, 1, '1']
$redis->close();
```

## 支持范围

- 连接：TCP、`tls://`（系统证书校验，TLS 1.2+）、Unix socket；AUTH 密码或 ACL 用户名/密码、SELECT、close、连接状态。
- 字符串与键：GET、SET（NX/XX/EX/PX/EXAT/PXAT/KEEPTTL/GET）、SETEX、PSETEX、SETNX、MGET、MSET、MSETNX、GETSET、GETDEL、DEL、UNLINK、EXISTS、TOUCH、过期时间、TTL/PTTL、PERSIST、TYPE、计数器、字符串范围、RENAME、KEYS。
- 哈希：HGET/HSET/HSETNX/HMGET/HMSET/HGETALL、HDEL、HEXISTS、HLEN、HKEYS/HVALS、HINCRBY/HINCRBYFLOAT。
- 列表：LPUSH/RPUSH/LPUSHX/RPUSHX、LPOP/RPOP、LLEN、LRANGE、LINDEX、LSET、LTRIM、LREM。
- 集合：SADD/SREM/SISMEMBER/SMEMBERS/SCARD/SPOP/SRANDMEMBER、SUNION/SINTER/SDIFF。
- 有序集合：ZADD、ZREM、ZCARD、ZSCORE、ZRANK/ZREVRANK、ZINCRBY、ZRANGE/ZREVRANGE（withScores 布尔参数）、ZCOUNT、按分数/排名删除范围。
- 脚本：EVAL/EVALSHA、SCRIPT；只对 KEYS 添加前缀，不修改 ARGV。
- 批处理：pipeline、MULTI/EXEC、DISCARD、WATCH/UNWATCH；使用同一物理连接维持事务与所选数据库。
- 其他：PING、ECHO、DBSIZE、TIME、INFO、FLUSHDB/FLUSHALL，以及 CLIENT 的 GETNAME/SETNAME/ID/INFO/LIST。
- rawCommand：发送原始参数，不添加前缀；拒绝切换协议、订阅模式、连接及事务状态的命令，以免破坏响应对应关系。
- 选项：OPT_PREFIX、OPT_READ_TIMEOUT，以及值为 NONE 的 OPT_SERIALIZER/OPT_COMPRESSION。
  `_prefix`、`_serialize`、`_unserialize`、`_pack`、`_unpack` 支持当前无序列化/压缩模式。
- 模块识别：extension_loaded、get_loaded_extensions、phpversion、ReflectionExtension。

方法名不区分大小写。PHP `lRem(key, value, count)` 的顺序会转换为 Redis 协议顺序。
HMGET 返回字段到值的映射；缺失值为 false，空字符串仍为空字符串。
普通服务端错误（如 WRONGTYPE）返回 false，并保存在 getLastError；认证、加载中、
网络和协议错误抛 RedisException。批处理会消费完整响应后再报告异常，避免下一条命令读到旧响应。

## 生命周期与限制

每个实例拥有独立连接和选项；方法串行访问该连接。显式 close 或垃圾回收会关闭连接。
网络操作使用当前 `runtime.RequestContext()`，请求取消会中断读写。
默认连接/读写超时为 60 秒，也可指定 connect/read_timeout 或 OPT_READ_TIMEOUT；
请求自身更短的截止时间优先。超时、断连和协议错误后丢弃连接，需要显式重新 connect，
不会自动重试可能已在服务器执行的写入。

尚未支持持久连接池/pconnect、自动重试与 backoff 选项、RedisCluster/Sentinel、
Pub/Sub、RESP3、阻塞命令的专用方法、SCAN 引用参数接口、PHP/JSON 序列化、压缩、
自定义 TLS context、构造函数连接选项、完整反射签名与命名参数契约。
批处理不能嵌套。未知选项或上述功能不会静默退化为其他行为。
命令返回值及转换以 phpredis 6.1 的相关接口为参照，但并非完整兼容实现。
单个响应或 pipeline 响应总量限制为 64 MiB，嵌套深度最多 128，批处理请求限制为 64 MiB。

Laravel 的默认 Redis 配置包含本模块未实现的重试/backoff 选项，因此不能据此宣称
Laravel Redis、缓存或队列已完整兼容。`std/laravel/framework/illuminate/redis`
仍是独立的框架模块入口，不会隐式加载此驱动。

## 验证

```powershell
go test -mod=mod ./std/php/redis ./std/php ./std/php/core

# 对支持 DB 0/1 的 standalone Redis 运行；仅删除随机前缀测试键，不执行 FLUSHDB
$env:ORIGAMI_REDIS_TEST_PORT='6379'
$env:ORIGAMI_REDIS_TEST_HOST='127.0.0.1'
go test -mod=mod ./std/php/redis -run TestRedisIntegration -v
Remove-Item Env:ORIGAMI_REDIS_TEST_PORT, Env:ORIGAMI_REDIS_TEST_HOST
```

Go 测试宿主显式调用 `redis.Load(vm)` 后运行 `tests/php/redis_driver_contract_test.php`
和 `tests/php/redis_driver_integration_test.php`，无需修改默认 CLI。
离线测试覆盖协议边界、二进制数据、请求取消、超时、命令参数映射、注册及实例隔离。
真实服务测试覆盖常用数据结构、Lua、pipeline、事务、WATCH、错误后继续读写与关闭。
若对独立临时服务设置 `ORIGAMI_REDIS_TEST_ACL=1`，还会创建并删除临时 ACL 用户以验证认证。
本次已在本机 Redis 8.2.1 临时实例上通过真实服务测试（含 ACL 和数据库隔离），实例已关闭。
本机 PHP 未加载 redis 扩展，尚未进行原生 phpredis 差分验收。

参考：[phpredis 6.1 接口](https://github.com/phpredis/phpredis/blob/6.1.0/redis.stub.php)、
[Redis RESP 协议](https://redis.io/docs/latest/develop/reference/protocol-spec/)。
