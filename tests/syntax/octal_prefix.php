<?php

namespace tests\syntax;

/**
 * PHP 8.1：0o / 0O 八进制字面量。
 */

if (0o17 !== 15) {
    Log::fatal('0o17 应为 15，实际 ' . 0o17);
}
if (0O10 !== 8) {
    Log::fatal('0O10 应为 8');
}

Log::info('syntax octal 0o 测试通过');
