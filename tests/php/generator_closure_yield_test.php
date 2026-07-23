<?php

namespace tests\php;

/**
 * Closure 内 yield（Symfony Table::buildTableRows 的 new TableRows(function(){ yield ...})）。
 */

$fn = function () {
    yield 1;
    yield 2;
};

$vals = [];
foreach ($fn() as $v) {
    $vals[] = $v;
}

if ($vals !== [1, 2]) {
    Log::fatal('closure yield failed: ' . var_export($vals, true));
}

Log::info('generator_closure_yield_test 测试通过');
