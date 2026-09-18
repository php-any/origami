<?php

namespace tests\tokenizer;

/**
 * tokenizer：token_get_all。
 */

$tokens = token_get_all('<?php echo 1;');
if (!is_array($tokens) || count($tokens) < 2) {
    Log::fatal('token_get_all 失败: ' . var_export($tokens, true));
}

Log::info('tokenizer 测试通过');
