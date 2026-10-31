# 通用网络脚本 API

本目录 API 不依赖具体框架。`zy.go` 注册 HTTP、WebSocket、TCP、gRPC 模块；自行装配 VM 时调用各包的 `Load(vm)`。

## TCP

```php
$server = new Net\Tcp\Server('127.0.0.1:9000');
$server->configure(['maxConnections'=>64, 'readTimeoutMs'=>5000,
    'writeTimeoutMs'=>5000, 'connectionTimeoutMs'=>30000]);
$server->serve(function ($connection) {
    try { $connection->write($connection->read(4096)); }
    finally { $connection->close(); }
});
```

`read(maxBytes)` 返回至多指定字节（1..16 MiB），EOF 返回空串，失败抛异常；不保证读满，分帧由应用完成。`write(bytes)` 写完全部内容或抛异常，返回写入字节数。每个连接使用独立 RequestVM。超限连接关闭，截止时间关闭 socket 并取消脚本上下文。未捕获的回调错误关闭当前连接。

## HTTP 和 WebSocket

`Net\Http\Server::configure(array)` 可设置 maxConnections、maxRequests、maxBodyBytes、maxHeaderBytes、readHeaderTimeoutMs、readTimeoutMs、writeTimeoutMs、idleTimeoutMs、requestTimeoutMs。请求体限制在脚本执行和超级全局变量读取前生效，包括 chunked 请求；超限返回 413，并发槽满返回 503。监听连接槽计入 keep-alive 和 hijacked 连接。未调用 configure 的已有服务器保持默认行为。

`Net\Websocket\Conn::configure(array)` 可设置 maxMessageBytes、readTimeoutMs、writeTimeoutMs。配置后，读消息限制在完整消息分配前执行，写消息也检查大小；连接随调用的请求上下文取消而关闭，close() 解除取消回调。只使用读取/写入超时，不能替代业务级状态检查。

## gRPC

```php
$schema = new Net\Grpc\Schema(['echo.proto'], ['/absolute/proto']);
$server = new Net\Grpc\Server('127.0.0.1:9001', $schema);
$server->configure(['maxCalls'=>32, 'maxRecvBytes'=>1048576]);
$server->reflection(true);
$server->serve(function ($call) {
    // Example echo contract: input and output have the same protobuf type.
    $call->respondJson($call->requestJson());
});
```

Schema 编译 proto、imports 和标准类型。`methodsJson()` 返回服务/方法名称、完整 path、input/output 类型及 clientStreaming/serverStreaming 标记。Server 当前只支持 unary 应用方法，streaming 契约明确报错；可选反射提供标准 v1/v1alpha 双向流。

Call 的 method() 返回完整 RPC 路径；requestJson()、respondJson(string) 遵循 ProtoJSON（普通消息对象、wrapper 标量、int64 字符串等）。fail(int code,string message) 发送 1..16 的标准 gRPC 错误。回调不响应、抛出未捕获异常或输出不符合契约时返回 INTERNAL。网络上使用 HTTP/2 Protobuf；启用通用 gzip codec，并限制解压后消息。每次调用使用独立 RequestVM。

可配置字段：maxConnections、maxConcurrentStreams（每连接）、maxCalls（进程级）、maxRecvBytes、maxSendBytes、maxHeaderBytes、requestTimeoutMs、connectionTimeoutMs（握手）、idleTimeoutMs、maxConnectionAgeMs。准入在解码前执行，传输准入配额持有至 RPC 完成或取消，反射流也受限制。执行中的 unary 回调另持有执行配额直至返回，忽略取消的原生扩展不能通过超时不断启动更多回调。期限从入站 header 开始，不等待消息解码后才计时。连接年龄包含 gRPC 自身的抖动和 requestTimeoutMs 结束宽限期。

以上所有大小单位为字节，超时单位为毫秒；选项必须是 1..1073741824 的整数，未知选项报错。默认值在各模块内可见。配置在启动/第一次 I/O 前完成，不支持并发重配置。监听为明文 TCP，没有 TLS 配置 API或统一信号停机编排；协作取消不能强行中止忽略 context 的原生扩展。

独立验证：

```powershell
go test -mod=readonly -timeout 30s ./std/net/tcp ./std/net/http ./std/net/websocket ./std/net/grpc ./std/net/internal/script
go run -mod=readonly ./zy.go tests/php/net_api_test.php
```

回调支持普通闭包与 `Closure::bindTo()` 绑定闭包，参数通过语言统一的形参绑定规则传入（包含类型检查、默认值、变长参数和 `func_get_args`）。
