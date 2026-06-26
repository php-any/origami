<?php

echo "=== 太空船运算符 <=> 测试 ===\n";

// 左边小于右边 => -1
$result = 1 <=> 2;
if ($result == -1) {
    Log::info("<=> 小于测试通过");
} else {
    Log::fatal("<=> 小于测试失败, 实际 {$result}");
}

// 左边等于右边 => 0
$result = 5 <=> 5;
if ($result === 0) {
    Log::info("<=> 等于测试通过");
} else {
    Log::fatal("<=> 等于测试失败, 实际 {$result}");
}

// 左边大于右边 => 1
$result = 10 <=> 3;
if ($result == 1) {
    Log::info("<=> 大于测试通过");
} else {
    Log::fatal("<=> 大于测试失败, 实际 {$result}");
}

// 字符串比较
$result = "a" <=> "b";
if ($result == -1) {
    Log::info("<=> 字符串小于测试通过");
} else {
    Log::fatal("<=> 字符串小于测试失败, 实际 {$result}");
}

$result = "abc" <=> "abc";
if ($result === 0) {
    Log::info("<=> 字符串相等测试通过");
} else {
    Log::fatal("<=> 字符串相等测试失败, 实际 {$result}");
}

$result = "z" <=> "a";
if ($result == 1) {
    Log::info("<=> 字符串大于测试通过");
} else {
    Log::fatal("<=> 字符串大于测试失败, 实际 {$result}");
}

// 用于排序
$unsorted = [5, 2, 8, 1, 9, 3];
$sorted = [];
// 手动冒泡排序使用 <=>
$arr = [5, 2, 8, 1, 9, 3];
$n = count($arr);
for ($i = 0; $i < $n - 1; $i++) {
    for ($j = 0; $j < $n - $i - 1; $j++) {
        if (($arr[$j] <=> $arr[$j + 1]) == 1) {
            $temp = $arr[$j];
            $arr[$j] = $arr[$j + 1];
            $arr[$j + 1] = $temp;
        }
    }
}
if ($arr[0] == 1 && $arr[1] == 2 && $arr[2] == 3 && $arr[3] == 5 && $arr[4] == 8 && $arr[5] == 9) {
    Log::info("<=> 排序测试通过");
} else {
    Log::fatal("<=> 排序测试失败");
}

// 浮点数比较
$result = 1.5 <=> 1.6;
if ($result == -1) {
    Log::info("<=> 浮点数比较测试通过");
} else {
    Log::fatal("<=> 浮点数比较测试失败, 实际 {$result}");
}

echo "=== 太空船运算符 <=> 测试完成 ===\n";
