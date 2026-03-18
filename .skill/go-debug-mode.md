# Go 调试模式规范

## 核心原则

**禁止使用 `go build` 进行调试**，直接使用 `go run` 执行和测试代码。

## 理由

1. **减少冗余步骤**：`go build` 生成二进制文件后再执行是多此一举
2. **提高效率**：`go run` 一步完成编译和执行
3. **保持上下文**：避免在构建和运行之间切换注意力
4. **符合 Go 哲学**：Go 工具链设计初衷就是快速迭代

## 正确做法

```bash
# ✅ 推荐：直接运行
go run origami.go laravel/artisan list

# ✅ 推荐：带参数运行
go run origami.go script.php --arg1=value

# ❌ 禁止：先编译再运行（多余步骤）
go build -o origami_bin origami.go
./origami_bin laravel/artisan list
```

## 例外情况

仅在以下场景允许使用 `go build`：

- 发布生产版本需要优化性能
- 需要分析二进制文件大小
- 特定平台的交叉编译测试
- 性能基准测试（避免编译时间干扰）

## 错误处理

当 `go run` 报错时：

1. 直接根据错误信息修复源代码
2. 重新执行 `go run` 验证修复
3. 不要先生成二进制文件再测试

## 工具集成

IDE 或编辑器配置应优先使用 `go run`：

```json
{
  "go.runCommand": "go run",
  "go.buildOnSave": false
}
```
