<?php

namespace tests\php;

/**
 * while 体内 yield（vlucas/phpdotenv Lexer::lex 的核心路径）。
 * 回归：未实现 WhileYieldControl 时生成器只会产出第一个元素。
 */

function gen_while_123() {
    $i = 0;
    while ($i < 3) {
        yield $i;
        $i++;
    }
}

$got = [];
foreach (gen_while_123() as $v) {
    $got[] = $v;
}
if ($got !== [0, 1, 2]) {
    Log::fatal('function while yield failed: ' . var_export($got, true));
}

$a = iterator_to_array(gen_while_123());
if (array_values($a) !== [0, 1, 2]) {
    Log::fatal('iterator_to_array while yield failed: ' . var_export($a, true));
}

// 贴近 Dotenv Lexer：isset 字符串下标 + while + yield
function dotenv_style_lex(string $content) {
    $offset = 0;
    while (isset($content[$offset])) {
        yield $content[$offset];
        $offset++;
    }
}

$chars = iterator_to_array(dotenv_style_lex('ab"'));
if (array_values($chars) !== ['a', 'b', '"']) {
    Log::fatal('dotenv-style while yield failed: ' . var_export($chars, true));
}

// 闭包
$fn = function () {
    $i = 0;
    while ($i < 2) {
        yield 'x' . $i;
        $i++;
    }
};
$got2 = iterator_to_array($fn());
if (array_values($got2) !== ['x0', 'x1']) {
    Log::fatal('closure while yield failed: ' . var_export($got2, true));
}

Log::info('generator_while_yield 测试通过');
