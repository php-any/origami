<?php
namespace tests\php;

$parts = ['APP_NAME', 'Laravel'];
[$name, $value] = $parts;
if ($value !== 'Laravel') {
    Log::fatal("destructure fail: ".var_export($value, true));
}

$fn = function () use ($value) {
    return $value;
};
if ($fn() !== 'Laravel') {
    Log::fatal("use capture fail: ".var_export($fn(), true));
}

// nested flatMap style
$outer = function (array $parts) {
    [$name, $value] = $parts;
    $inner = function () use ($value) {
        return $value;
    };
    return $inner();
};
if ($outer(['APP_NAME', 'Laravel']) !== 'Laravel') {
    Log::fatal('nested use fail');
}

Log::info('destructure_use_closure 测试通过');
