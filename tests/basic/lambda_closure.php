<?php

echo "=== 闭包和 Lambda 高级测试 ===\n";

// 闭包作为参数传递
function applyFunc($value, $callback) {
    return $callback($value);
}

$result = applyFunc(5, function($x) { return $x * $x; });
if ($result == 25) {
    Log::info("闭包作为参数测试通过");
} else {
    Log::fatal("闭包作为参数测试失败, 实际 {$result}");
}

// 闭包 use 值捕获
$multiplier = 3;
$closure = function($x) use ($multiplier) {
    return $x * $multiplier;
};
if ($closure(4) == 12) {
    Log::info("闭包 use 值捕获测试通过");
} else {
    Log::fatal("闭包 use 值捕获测试失败");
}

// 闭包 use 引用捕获
$counter = 0;
$increment = function() use (&$counter) {
    $counter++;
    return $counter;
};
$increment();
$increment();
$increment();
if ($counter == 3) {
    Log::info("闭包 use 引用捕获测试通过");
} else {
    Log::fatal("闭包 use 引用捕获测试失败, 实际 {$counter}");
}

// 闭包工厂
function makeMultiplier($factor) {
    return function($x) use ($factor) {
        return $x * $factor;
    };
}
$double = makeMultiplier(2);
$triple = makeMultiplier(3);
if ($double(5) == 10 && $triple(5) == 15) {
    Log::info("闭包工厂测试通过");
} else {
    Log::fatal("闭包工厂测试失败");
}

// 闭包数组操作
$numbers = [1, 2, 3, 4, 5];
$squared = [];
for ($n in $numbers) {
    $squared[] = (function($x) { return $x * $x; })($n);
}
if ($squared[0] == 1 && $squared[1] == 4 && $squared[2] == 9 && $squared[3] == 16 && $squared[4] == 25) {
    Log::info("闭包内联调用测试通过");
} else {
    Log::fatal("闭包内联调用测试失败");
}

// 闭包返回闭包
function makeGreeter($greeting) {
    return function($name) use ($greeting) {
        return $greeting . ", " . $name . "!";
    };
}
$hello = makeGreeter("Hello");
if ($hello("World") == "Hello, World!") {
    Log::info("闭包返回闭包测试通过");
} else {
    Log::fatal("闭包返回闭包测试失败");
}

// 闭包作为数组元素
$ops = [
    function($x) { return $x + 1; },
    function($x) { return $x * 2; },
    function($x) { return $x - 3; },
];
$val = 5;
for ($op in $ops) {
    $val = $op($val);
}
// ((5+1)*2)-3 = 9
if ($val == 9) {
    Log::info("闭包作为数组元素测试通过");
} else {
    Log::fatal("闭包作为数组元素测试失败, 实际 {$val}");
}

// 静态闭包
$staticVal = 10;
$staticClosure = static function($x) use ($staticVal) {
    return $x + $staticVal;
};
if ($staticClosure(5) == 15) {
    Log::info("静态闭包测试通过");
} else {
    Log::fatal("静态闭包测试失败");
}

echo "=== 闭包和 Lambda 高级测试完成 ===\n";
