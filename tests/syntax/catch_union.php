<?php

namespace tests\syntax;

/**
 * PHP 8.0：catch (A|B $e) 联合类型。
 */

$hit = '';
try {
    throw new \InvalidArgumentException('u');
} catch (\RuntimeException|\InvalidArgumentException $e) {
    $hit = $e->getMessage();
}

if ($hit !== 'u') {
    Log::fatal('联合 catch 失败: ' . $hit);
}

Log::info('syntax catch union 测试通过');
