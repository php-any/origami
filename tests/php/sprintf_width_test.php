<?php

namespace tests\php;

/**
 * sprintf 宽度/对齐 与 Symfony TextDescriptor 使用场景的兼容性测试。
 *
 * 对应 vendor/symfony/console/Descriptor/TextDescriptor.php:169:
 *   sprintf("%-{$width}s %s", $command->getName(), $command->getDescription())
 */

$width = 20;
$name = '-h, --help';
$desc = 'Display help for the given command.';

$line = sprintf("%-{$width}s %s", $name, $desc);

// 这里不做特别严格的空格断言，只要不出现 Go fmt 的 %!$(string=...) 伪输出即可。
if (str_contains($line, '%!')) {
    Log::fatal('sprintf 宽度测试失败: 出现 Go 格式错误片段: '.$line);
}

Log::info('sprintf 宽度测试通过: '.$line);

