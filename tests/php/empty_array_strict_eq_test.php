<?php

namespace tests\php;

$a = [];
if (!($a === [])) {
    Log::fatal('[] === [] 应为 true');
}

$b = array_values([]);
if (!($b === [])) {
    Log::fatal('array_values([]) === [] 应为 true');
}
Log::info('empty_array_strict_eq 测试通过');
