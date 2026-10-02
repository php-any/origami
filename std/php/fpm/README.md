# PHP 请求宿主适配（`std/php/fpm`）

语言层只有 `runtime.VM` 和 `runtime.RequestVM`。本包不实现独立的 VM：`fpm.RequestVM` 是 `runtime.RequestVM` 的类型别名，`fpm.New` 只创建请求 VM 并绑定输出。Laravel、经典 PHP 入口和解析预览使用同一个 RequestVM 实现。

VM 提供标准库、启动期声明、命名空间路径和共享程序缓存。RequestVM 持有请求新增声明、常量、全局变量、会话、include 状态、函数 static 局部变量、错误/异常处理器、shutdown 队列、调用栈和输出。缓存程序的请求声明不会登记到共享 VM；启动期已执行文件集合以不可变快照继承，后续 include 写入请求自己的映射。

```go
p := parser.NewParser()
vm := runtime.NewVM(p).(*runtime.VM)
php.Load(vm)

request := fpm.New(vm, func(s string) { _, _ = io.WriteString(w, s) })
request.BindHTTP(r, w)
_, control := request.LoadAndRun("index.php")
if control != nil {
    _, control = request.HandleUnhandledException(control)
}
// 宿主负责处理剩余控制流、HTTP 状态和已提交响应。
request.RunShutdownCallbacks()
```

宿主必须在失败路径也完成 shutdown 和输出收尾，HTTP 客户端须设置超时。`LoadAndRunFresh` 跳过解析缓存，用于热重载；`RunCompiledFile` 在 RequestVM 上执行预编译程序，不回到 VM 的执行上下文。

Laravel 的容器/Facade/原生服务适配仍需提供正确的请求对象状态；统一 VM 类型不等于已证明任意捕获对象图或原生 singleton 的隔离。

验证：`go test -race ./runtime ./std/php/fpm`。回归覆盖共享程序的请求声明、全局变量和常量、编译文件执行、include 状态、HTTP/输出绑定、处理器栈及并发请求。
