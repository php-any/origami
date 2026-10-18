# 语言测试

本目录保留 Origami 语言特性、PHP 语法与运行时语义、PHP 内置函数及标准扩展的回归测试，以及必要的测试夹具。

从框架问题缩减出的独立语言用例可以保留；测试本身不应依赖 Laravel、Symfony、Livewire、Filament 或示例应用的 vendor、配置和服务。框架与业务模块测试应放在对应模块或示例的测试目录中。

从仓库根目录运行单个 PHP 回归：

```powershell
go run ./zy.go tests/php/arrayaccess_property_append_test.php
```

`run_tests.php` 是旧式 PHP 脚本批量入口，不覆盖 `.zy`、`.phpt` 和 Go 测试。核心临时文件使用系统临时目录，避免依赖示例应用目录。
