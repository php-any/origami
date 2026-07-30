<?php

namespace tests\php;

/**
 * getcwd 基本行为
 */

$cwd = getcwd();
if ($cwd === false || !is_string($cwd) || $cwd === '') {
    Log::fatal('getcwd 应返回非空字符串');
}
if (!is_dir($cwd)) {
    Log::fatal('getcwd 返回路径不是目录: '.$cwd);
}

Log::info('getcwd 测试通过: '.$cwd);
