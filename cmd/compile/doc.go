// Package compile 将 PHP 源码预编译为 Go AST 字面量。
//
// # 生成流程
//
//	PHP 解析 → Generator.Emit 反射转译 → go/format 格式化 → 写入输出目录
//
// # 新增 AST 节点时如何补充生成规则
//
// Emit 按以下优先级处理节点，命中即停止；全部失败则返回 EmitError 终止编译：
//
//  1. special_handlers.go 中的特殊处理器（specialHandlers 注册表）
//  2. emit_data.go 中的 data 标量 Value（IntValue、StringValue 等）
//  3. 反射结构体字面量（emitStructLiteral）自动处理。
//
// ## 情形 A：节点字段均可导出，嵌入的 *Node 带 pp:"-" 标签
//
// 无需任何改动。Emit 会自动通过反射遍历字段，生成 &node.Xxx{Node: node.NewNode(from), ...}。
// 确保嵌入的 *Node 带 pp:"-" 标签，所有字段已导出。
//
// ## 情形 B：含未导出字段、运行时引用或需转换逻辑
//
// 在 special_handlers.go 中：
//
//  1. 实现 emitXxx(g *Generator, v data.GetValue) error
//  2. 在 specialHandlers 注册 reflect.TypeOf((*node.Xxx)(nil))
//
// 典型场景：CallExpression（Fun 不可序列化 → NewCallTodo）、ClassStatement（Methods map）、
// VarFastAssign（未导出 op → node.NewVarFastAssignCompiled）。
//
// 若未导出字段仅为 primitive/string，可在特殊处理器中用 reflect 读取，或于 node 包
// 添加 Compile 专用构造函数 / getter。
//
// ## 情形 C：新增 data 标量或注解类型
//
// 标量 Value：在 emit_data.go 的 dataValueEmitters 注册发射函数。
// 框架注解：在 special_handlers.go 的 emitClassAnnotation 增加 case，
// 并在 std/.../compile_bootstrap.go 提供 CompiledXxxValue 工厂（如需要）。
//
// ## 验证
//
// 运行 zy compile <dir> -o <out> 触发编译；不支持的节点会报错并附带文件路径与行列号。
// 生成产物经 go/format 格式化，可直接 go build。
package compile
