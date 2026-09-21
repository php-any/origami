<?php

namespace tests\php;

/**
 * htmlspecialchars 与空字符串边界（Illuminate e() 底层），避免加速层改坏转义。
 */

if (htmlspecialchars('') !== '') {
    \Log::fatal('empty htmlspecialchars');
}
if (htmlspecialchars('fullWidth') !== 'fullWidth') {
    \Log::fatal('plain htmlspecialchars');
}
if (htmlspecialchars('<x>') !== '&lt;x&gt;') {
    \Log::fatal('lt htmlspecialchars: '.var_export(htmlspecialchars('<x>'), true));
}

\Log::info('htmlspecialchars illuminate 边界测试通过');
