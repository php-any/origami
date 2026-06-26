<?php

// 测试 CLI 注解

use Cli\Annotation\CliApplication;
use Cli\Annotation\Command;

// 定义一个测试命令
#[Command(name: "test", description: "测试命令")]
class TestCommand {
    public function execute(): void {
        echo "Test command executed!\n";
    }
}

// 定义另一个测试命令
#[Command(name: "hello", description: "问候命令")]
class HelloCommand {
    public function execute(): void {
        global $argv;

        $name = "World";
        if (count($argv) > 0) {
            $name = $argv[0];
        }

        echo "Hello, " . $name . "!\n";
    }
}

// CLI 应用入口
#[CliApplication(name: "TestCLI", version: "1.0.0")]
class TestCliApp {
    public static function boot(): void {
        echo "Test CLI Application started\n";
    }

    public static function exit(): void {
        echo "Test CLI Application shutting down\n";
    }
}

echo "CLI annotation test completed\n";
