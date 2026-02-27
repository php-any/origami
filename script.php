<?php

/**
 * 最小复现场景：数组按“引用”共享被提前消费。
 *
 * 在 PHP 官方 CLI 中，下面代码的语义是：
 *
 *   $tokens = $argv;       // 拷贝一份数组
 *   array_shift($tokens);  // 只修改 $tokens，不影响 $argv
 *
 * 但在当前 Origami 实现中，数组赋值更像“共享底层存储”，
 * 导致对 $tokens 的 array_shift 会把 $argv 一并改掉。
 *
 * 运行对比：
 *   1) php script.php hello 33
 *   2) go run ./origami.go ./script.php hello 33
 */

// 用真实的 CLI 参数来演示（更贴近 ArgvInput/$_SERVER['argv'] 场景）
echo "==== 原始 \$argv ====\n";
var_dump($argv);

// 这里模拟诸如 ArgvInput::getParameterOption 中的 `$tokens = $this->tokens` 赋值
$tokens = $argv;

// 消费一部分 tokens（比如解析完脚本名/命令名后向前推进）
array_shift($tokens);

echo "\n==== array_shift(\$tokens) 之后 ====\n";
echo "\$argv 现在是：\n";
var_dump($argv);

echo "\n\$tokens 是：\n";
var_dump($tokens);

echo "\n==== 预期（PHP CLI） vs 实际（Origami 当前实现）说明 ====\n";
echo "- 在 PHP 中：\$argv 仍然包含完整参数（不会被上面的 array_shift 改动）。\n";
echo "- 在 Origami 中：若你看到 \$argv 与 \$tokens 同时被截短，就说明数组按引用共享被提前消费了。\n";

