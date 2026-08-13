<?php

namespace tests\php;

/**
 * PHP 8+：字符串连接 . 优先级低于 +、- 与位移。
 */

if (("Number: " . 1 + 2) !== "Number: 3") {
    Log::fatal('. 低于 +: 期望 Number: 3, 实际 ' . var_export("Number: " . 1 + 2, true));
}

if ((1 + 2 . "x") !== "3x") {
    Log::fatal('+ 高于 .: 期望 3x');
}

if (("a" . 1 << 2) !== "a4") {
    Log::fatal('. 低于 <<: 期望 a4');
}

if ((2 * 3 . 4) !== "64") {
    Log::fatal('* 高于 .: 期望 64');
}

Log::info('concat_precedence_php8 测试通过');
