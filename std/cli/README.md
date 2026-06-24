# CLI 注解模块

CLI 注解模块为 Origami 语言提供了命令行应用开发支持，类似于 Symfony Console 的注解方式。

## 功能特性

- `#[CliApplication]` - CLI 应用入口注解
- `#[Command]` - 命令定义注解
- 自动命令注册和分发
- 命令行参数解析
- 帮助信息生成

## 使用示例

### 1. 定义命令

```php
#[Command(name: "greet", description: "向某人打招呼")]
class GreetCommand {
    public function execute(): void {
        global $argv;

        $name = "World";
        if (count($argv) > 0) {
            $name = $argv[0];
        }

        echo "Hello, " . $name . "!\n";
    }
}
```

### 2. 定义 CLI 应用

```php
#[CliApplication(name: "MyCLI", version: "1.0.0")]
class App {
    public static function boot(): void {
        // 应用启动时的初始化逻辑
    }

    public static function exit(): void {
        // 应用退出时的清理逻辑
    }
}
```

### 3. 运行应用

```bash
./origami my_cli_app.zy greet World
```

## 注解参数

### CliApplication

| 参数 | 类型 | 默认值 | 描述 |
|------|------|--------|------|
| name | string | "CLI" | 应用名称 |
| version | string | "1.0.0" | 应用版本 |

### Command

| 参数 | 类型 | 默认值 | 描述 |
|------|------|--------|------|
| name | string | "" | 命令名称 |
| description | string | "" | 命令描述 |

## 内置命令

- `--help`, `-h` - 显示帮助信息
- `--version`, `-v` - 显示版本信息

## API 参考

### CliRuntime

CLI 运行时类，负责命令分发和执行。

```go
runtime := cli.NewCliRuntime("MyCLI", "1.0.0")
runtime.Run(ctx)
```

### OptionParser

选项解析器，用于解析命令行选项。

```go
parser := cli.NewOptionParser()
parser.AddOption(cli.Option{
    Name:        "name",
    ShortName:   "n",
    Description: "用户名称",
    Required:    true,
    HasValue:    true,
})

options, args, err := parser.Parse(os.Args[1:])
```

### CommandBase

命令基类，提供通用功能。

```go
type MyCommand struct {
    *cli.CommandBase
}

func NewMyCommand() *MyCommand {
    cmd := &MyCommand{
        CommandBase: cli.NewCommandBase("mycommand", "我的命令"),
    }
    cmd.AddOption(cli.Option{
        Name:        "verbose",
        ShortName:   "v",
        Description: "详细输出",
    })
    return cmd
}
```

## 示例代码

查看 `examples/cli_app.php` 文件获取完整示例。
