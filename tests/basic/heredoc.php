<?php

echo "=== Heredoc 和 Nowdoc 字符串测试 ===\n";

// 基本 heredoc
$str = <<<EOT
Hello World
EOT;
if ($str == "Hello World") {
    Log::info("基本 heredoc 测试通过");
} else {
    Log::fatal("基本 heredoc 测试失败, 实际 '{$str}'");
}

// Heredoc 多行
$multi = <<<TEXT
Line 1
Line 2
Line 3
TEXT;
if ($multi == "Line 1\nLine 2\nLine 3") {
    Log::info("heredoc 多行测试通过");
} else {
    Log::fatal("heredoc 多行测试失败");
}

// Heredoc 中使用变量插值
$name = "Origami";
$interpolated = <<<MSG
Hello, $name!
MSG;
if ($interpolated == "Hello, Origami!") {
    Log::info("heredoc 变量插值测试通过");
} else {
    Log::fatal("heredoc 变量插值测试失败, 实际 '{$interpolated}'");
}

// Heredoc 用不同定界符
$html = <<<HTML
<div>Test</div>
HTML;
if ($html == "<div>Test</div>") {
    Log::info("heredoc HTML 定界符测试通过");
} else {
    Log::fatal("heredoc HTML 定界符测试失败");
}

// Heredoc 空内容
$empty = <<<EMPTY

EMPTY;
if ($empty == "") {
    Log::info("heredoc 空内容测试通过");
} else {
    Log::fatal("heredoc 空内容测试失败");
}

echo "=== Heredoc 和 Nowdoc 字符串测试完成 ===\n";
