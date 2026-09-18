<?php

namespace tests\tokenizer;

/**
 * tokenizer：PHP 8 名称 token（T_NAME_FULLY_QUALIFIED 等常量存在）。
 */

if (!defined('T_STRING')) {
    Log::fatal('T_STRING 未定义');
}
$tokens = token_get_all('<?php echo \\Foo\\Bar;');
if (!is_array($tokens) || count($tokens) < 2) {
    Log::fatal('token_get_all 命名空间名失败');
}

Log::info('tokenizer name tokens 测试通过');
