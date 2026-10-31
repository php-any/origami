# MongoDB 驱动

MongoDB 是可选模块，`std/php.Load` 不会自动注册。底层使用官方 `go.mongodb.org/mongo-driver/v2`。
这是 PHP `ext-mongodb` API 的初始子集，模块版本为 `0.1.0-origami`，
不宣称完整兼容任一版本的官方 PHP 扩展，也不保证 Composer `mongodb/mongodb`
或 Laravel MongoDB 集成包可以直接运行。

## 使用

在需要 MongoDB 的宿主 `main.go` 中导入模块，并在创建 VM 时主动注册：

```go
import (
    "github.com/php-any/origami/data"
    "github.com/php-any/origami/std/php"
    "github.com/php-any/origami/std/php/mongodb"
)

func loadPHP(vm data.VM) {
    php.Load(vm)
    mongodb.Load(vm)
}
```

将 `mongodb.Load(vm)` 加入宿主实际使用的 VM 加载链。未注册时，MongoDB 类、函数、
接口和版本常量均不可用，`extension_loaded('mongodb')` 返回 false。
默认 CLI `zy.go` 和框架示例不会主动启用此模块。

注册后可运行以下 PHP 代码：

```php
<?php
$manager = new MongoDB\Driver\Manager('mongodb://127.0.0.1:27017/', [
    'serverSelectionTimeoutMS' => 3000,
    'timeoutMS' => 5000,
]);
$bulk = new MongoDB\Driver\BulkWrite();
$id = $bulk->insert(['name' => 'Origami', 'active' => true]);
$result = $manager->executeBulkWrite('demo.items', $bulk);
echo $result->getInsertedCount();

$query = new MongoDB\Driver\Query(['_id' => $id]);
$cursor = $manager->executeQuery('demo.items', $query);
$cursor->setTypeMap(['root' => 'array', 'document' => 'array']);
foreach ($cursor as $document) {
    echo $document['name'];
}
```

Manager 构造只校验配置，首次执行操作时才创建客户端。URI 中的连接、认证、TLS、
副本集等选项由官方 Go 驱动处理。Manager 的 `uriOptions` 数组支持常用超时、
池大小、认证源、读偏好、重试、TLS 和 write concern 选项；非空 `driverOptions`
尚未支持。每个 Manager 持有自己的连接池，游标持有 Manager 状态，避免客户端被提前回收。
操作使用 Origami 当前请求的取消上下文；垃圾回收时关闭游标和客户端，清理超时为 5 秒。

## 已实现

- `Manager::executeQuery / executeCommand / executeBulkWrite`。
- `Query`：filter、projection、sort、limit、skip、batchSize、hint、comment、allowDiskUse、noCursorTimeout。
- `Command`：普通命令返回单条结果游标；find、aggregate、listCollections、listIndexes 处理服务端游标。
- `BulkWrite`：insert、操作符 update、replacement update、delete、count；ordered、upsert、multi、delete limit、bypassDocumentValidation。
  入队时复制 BSON 文档；生成 `_id` 不修改调用方数组；同一批次不能重复执行。
- `WriteConcern`：w、journal；wtimeout 目前仅接受 0，操作超时可使用 URI 的 timeoutMS。
- `WriteResult`：插入、匹配、修改、删除、upsert 计数和 upsert ID，isAcknowledged。
- `Cursor`：getIterator、toArray、setTypeMap、isDead，按需 getMore；推进后不能回退。
- BSON：ObjectId、UTCDateTime（毫秒整数）、Binary、Regex、Decimal128、Timestamp、MinKey、MaxKey，及 JSON 序列化。
- `MongoDB\BSON\fromPHP / toPHP / fromJSON / toCanonicalExtendedJSON / toRelaxedExtendedJSON`。
- 连接、超时、命令及批量写入异常；保留服务端错误码和消息，BulkWriteException 可读取部分 WriteResult。
- `extension_loaded`、`get_loaded_extensions`、`phpversion`、ReflectionExtension 的模块识别。

BSON 文档保持字段顺序。顶层数组编码为文档；嵌套连续整数键数组编码为 BSON 数组，
稀疏或关联数组编码为文档。默认文档解码为 stdClass，数组仍为 PHP 数组。
type map 仅支持 root/document 的 array、object、stdClass，以及 array 的 array。
递归或超过 100 层的输入报错，不吞掉转换错误。

## 边界

尚未实现 Session/事务、监控订阅、Server API、ReadPreference/ReadConcern 对象、
ClientEncryption、客户端级 BulkWriteCommand、自定义 BSON Serializable/Unserializable、
自定义类 type map、更新 pipeline、collation/arrayFilters、完整的错误详情对象、
UTCDateTime::toDateTime，以及完整反射/序列化契约。未知数组选项明确报错。
不要将这一子集用于依赖上述能力的应用而不先补齐并验证。

## 验证

```powershell
go test -mod=mod ./std/php/mongodb ./std/php ./std/php/core
# 测试宿主显式注册 MongoDB 后执行 tests/php/mongodb_bson_driver_test.php
go test -mod=mod ./std/php/mongodb -run 'TestMongoDBExplicitRegistration|TestMongoDBPHPRegression' -v

# 可选的真实服务验收；仅创建/删除 origami_test 下唯一命名的测试集合
$env:ORIGAMI_MONGODB_URI='mongodb://127.0.0.1:27017/?serverSelectionTimeoutMS=3000&timeoutMS=5000'
go test -mod=mod ./std/php/mongodb -run TestMongoDBIntegration -v
Remove-Item Env:ORIGAMI_MONGODB_URI
```

包测试含本地 MongoDB 协议测试服务，用于验证真实驱动的网络请求、getMore 和错误传播；
它不是 MongoDB 服务端实现，不能代替可选的真实服务测试。
当前环境的 PHP 8.1.34 未加载 mongodb 扩展，无法直接进行原生 PHP 差分验证。

行为参考：[PHP MongoDB 驱动手册](https://www.php.net/manual/en/book.mongodb.php)、
[官方 Go 驱动文档](https://www.mongodb.com/docs/drivers/go/current/)。
