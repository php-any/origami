<?php

namespace tests\php;

/**
 * 测试 PHP 风格的 (object) 转换：
 * - 数组 (含数值键) 转为对象，其键作为属性名
 * - 标量值被包装到带有 scalar 属性的对象中
 */

// 1. 数组 -> 对象
$arr = ['name' => 'Taylor', 'age' => 36];
$obj = (object) $arr;

if (!isset($obj->name) || $obj->name !== 'Taylor') {
    Log::fatal('object_cast_test: 数组到对象的转换 name 属性不符合预期');
}

if (!isset($obj->age) || $obj->age !== 36) {
    Log::fatal('object_cast_test: 数组到对象的转换 age 属性不符合预期');
}

// 2. 数值键数组 -> 对象（键应变为 "0", "1", ...）
$list = ['a', 'b'];
$obj2 = (object) $list;

if (!isset($obj2->{'0'}) || $obj2->{'0'} !== 'a') {
    Log::fatal('object_cast_test: 数值键数组到对象的转换 0 属性不符合预期');
}

if (!isset($obj2->{'1'}) || $obj2->{'1'} !== 'b') {
    Log::fatal('object_cast_test: 数值键数组到对象的转换 1 属性不符合预期');
}

// 3. 标量 -> 对象（scalar 属性）
$scalarObj = (object) 42;
if (!isset($scalarObj->scalar) || $scalarObj->scalar !== 42) {
    Log::fatal('object_cast_test: 标量到对象的转换 scalar 属性不符合预期');
}

Log::info('object_cast_test: (object) 类型转换行为测试通过');

