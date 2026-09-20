<?php

namespace tests\php;

/**
 * !== / === 在一侧为 null 时不得 panic，应对齐 PHP 严格比较。
 */

$a = null;
$b = 'fi-icon-btn';
if (!($a !== $b)) {
    \Log::fatal('null !== string 应为 true');
}
if ($a === $b) {
    \Log::fatal('null === string 应为 false');
}
if (!($a === null)) {
    \Log::fatal('null === null 应为 true');
}
if (false !== 0) {
    // ok
} else {
    \Log::fatal('false !== 0 应为 true');
}
if (!('' !== null)) {
    \Log::fatal('空串 !== null 应为 true');
}

\Log::info('严格比较 null 测试通过');
