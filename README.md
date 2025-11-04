# Origami-lang

Origami-lang is an innovative hybrid scripting language that deeply integrates PHP's rapid development capabilities with Go's efficient concurrency model. It also incorporates some Java and TypeScript conventions.

> [中文文档](README_CN.md) | [English Documentation](README.md)

## ⚠️ Current Status

The codebase has not been optimized yet, and performance is not optimized.
Please use it as a tool, do not use it in production environments.

## 🚀 Core Features

### 🎯 Go Reflection Integration

- **Easy Registration**: Register Go functions and structs to the script domain with zero configuration
- **Automatic Type Conversion**: Seamless integration between Go and script types
- **Named Parameters**: Support for constructor named parameters

### 🎨 Syntax Fusion

- **PHP Compatibility**: Supports most PHP syntax patterns
- **Go Concurrency**: `spawn` keyword for coroutines
- **Type System**: Type declarations and nullable types
- **Generic Classes**: Support for generic syntax `class DB<T>`

### 💡 Special Features

- **HTML Embedding**: Direct HTML code blocks
- **Duck Typing**: `like` keyword for structural matching
- **Chinese Keywords**: Support for Chinese programming keywords
- **Functional Programming**: Rich array methods (`map`, `filter`, `reduce`, etc.)

## 🚀 Quick Start

```bash
git clone https://github.com/php-any/origami.git
cd origami
go build -o origami .
./origami script.zy
```

## 📚 Documentation

For detailed documentation, please visit the [Documentation Center](docs/README.md):

- **[Quick Start](docs/quickstart.md)** - Get started in 5 minutes
- **[Language Reference](docs/)** - Complete language documentation
  - [Syntax](docs/syntax.md) - Language syntax
  - [Data Types](docs/data-types.md) - Supported data types
  - [Functions](docs/functions.md) - Function definitions
  - [Classes](docs/classes.md) - Object-oriented programming
  - [Array Methods](docs/array_methods.md) - Array operations
- **[Go Integration](docs/go-integration.md)** - Integrate Go functions and structs
- **[Database Module](docs/database.md)** - Complete ORM documentation
- **[Standard Library](docs/std/)** - Built-in library reference

## 📝 Examples

See the [examples directory](examples/) for complete examples, or visit the [test cases](tests/) for more usage patterns.

## 💬 Discussion Group

<img src="https://github.com/php-any/origami/blob/main/qrcode_1753692981069.jpg" alt="Origami Discussion Group QR Code" width="200" />

## 📄 License

MIT License
