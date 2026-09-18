<?php

namespace tests\php;

/**
 * array_search(..., true) 对闭包必须按身份匹配，不能把所有 Closure 当成同一个。
 * Livewire EventBus::off() 依赖此行为。
 */

$a = function () {
    return 1;
};
$b = function () {
    return 2;
};
$list = [$a, $b];

$idx = array_search($b, $list, true);
if ($idx !== 1) {
    Log::fatal('array_search 严格模式应找到第二个闭包, got=' . var_export($idx, true));
}

$idx2 = array_search($a, $list, true);
if ($idx2 !== 0) {
    Log::fatal('array_search 严格模式应找到第一个闭包, got=' . var_export($idx2, true));
}

$c = function () {
    return 3;
};
$idx3 = array_search($c, $list, true);
if ($idx3 !== false) {
    Log::fatal('不存在的闭包应返回 false, got=' . var_export($idx3, true));
}

// 标量严格
if (array_search(1, ['1', 1], true) !== 1) {
    Log::fatal('严格模式下 1 不应匹配字符串 "1"');
}

Log::info('array_search 闭包身份测试通过');
