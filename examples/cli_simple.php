<?php

// CLI 简单示例 - 展示命令行参数处理

echo "=== CLI Simple Example ===\n\n";

// 获取命令行参数
global $argv;

echo "Command line arguments:\n";
echo "  argc: " . count($argv) . "\n";
echo "  argv: [";
for ($i = 0; $i < count($argv); $i++) {
    if ($i > 0) echo ", ";
    echo $argv[$i];
}
echo "]\n\n";

// 处理命令
if (count($argv) > 1) {
    $command = $argv[1];

    switch ($command) {
        case "greet":
            $name = count($argv) > 2 ? $argv[2] : "World";
            echo "Hello, " . $name . "!\n";
            break;

        case "version":
            echo "CLI Example version 1.0.0\n";
            break;

        case "help":
            echo "Available commands:\n";
            echo "  greet [name]  - Greet someone\n";
            echo "  version       - Show version\n";
            echo "  help          - Show this help\n";
            break;

        default:
            echo "Unknown command: " . $command . "\n";
            echo "Use 'help' to see available commands\n";
    }
} else {
    echo "Usage: php cli_simple.php <command> [arguments]\n";
    echo "Use 'help' to see available commands\n";
}
