// Package perfmon 是 Origami 的可选性能监控。
//
// 默认（无 build tag）全部为空实现，由编译器内联掉，不影响上线。
// 诊断时加上 -tags origamidebug 才会编入请求计时、autoload/解析计数和 pprof。
//
//	go run -mod=mod -tags origamidebug . serve --port=18086
//
// 环境变量（仅 debug 构建生效）：
//
//	ORIGAMI_CPUPROFILE=<path>   CPU profile，进程退出时落盘
//	ORIGAMI_MEMPROFILE=<path>   退出时堆快照
//	ORIGAMI_TRACE=<path>        execution trace
//	ORIGAMI_PPROF_ADDR=<addr>   另开 HTTP pprof（如 127.0.0.1:6060）
package perfmon
