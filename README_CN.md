# 折言(origami-lang)

折言(origami-lang) 是一门创新性的融合型脚本语言，深度结合 PHP 的快速开发基因与 Go 的高效并发模型。同时还有部分 Java、TypeScript 习惯引入。

> [中文文档](README_CN.md) | [English Documentation](README.md)

## ⚠️ 当前状态

当前未对代码分支进行任何优化，性能尚未优化。
请作为一个工具使用，请勿用于生产环境。

## 🚀 核心特征

### 🎯 Go 反射集成

- **便捷注册**: 一键将 Go 函数和结构体注册到脚本域，零配置
- **自动类型转换**: Go 和脚本类型之间的无缝集成
- **命名参数**: 支持构造函数命名参数

### 🎨 语法融合

- **PHP 兼容**: 支持大部分 PHP 语法模式
- **Go 并发**: `spawn` 关键字启动协程
- **类型系统**: 类型声明和可空类型
- **泛型类**: 支持泛型语法 `class DB<T>`

### 💡 特殊特性

- **HTML 内嵌**: 直接内嵌 HTML 代码块
- **鸭子类型**: `like` 关键字进行结构匹配
- **中文关键字**: 支持中文编程关键字
- **函数式编程**: 丰富的数组方法（`map`、`filter`、`reduce` 等）

## 🚀 快速开始

```bash
git clone https://github.com/php-any/origami.git
cd origami
go build -o origami .
./origami script.zy
```

## 📚 文档

详细文档请访问 [文档中心](docs/README.md)：

- **[快速开始](docs/quickstart.md)** - 5 分钟快速上手
- **[语言参考](docs/)** - 完整的语言文档
  - [语法](docs/syntax.md) - 语言语法
  - [数据类型](docs/data-types.md) - 支持的数据类型
  - [函数](docs/functions.md) - 函数定义
  - [类和对象](docs/classes.md) - 面向对象编程
  - [数组方法](docs/array_methods.md) - 数组操作
- **[Go 集成](docs/go-integration.md)** - 集成 Go 函数和结构体
- **[数据库模块](docs/database.md)** - 完整的 ORM 文档
- **[标准库](docs/std/)** - 内置库参考

## 📝 示例

查看 [示例目录](examples/) 了解完整示例，或访问 [测试用例](tests/) 查看更多使用模式。

## 💬 讨论群

<img src="https://github.com/php-any/origami/blob/main/qrcode_1753692981069.jpg" alt="折言讨论群二维码" width="200" />

## 📄 许可证

MIT 许可证
