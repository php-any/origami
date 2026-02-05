<?php

use Log;

echo "=== ->{表达式} 动态属性/方法 测试 ===\n";

class DynProp {
    public string $foo = 'foo';
    public string $bar = 'bar';

    public function getName(): string {
        return 'foo';
    }

    public function hello(string $name): string {
        return "hello-".$name;
    }
}

$obj = new DynProp();

// 1. $obj->{$name}
$prop = 'foo';
$v1 = $obj->{$prop};
if ($v1 === 'foo') {
    Log::info('obj->{$prop} ok');
} else {
    Log::fatal('obj->{$prop} fail, got '.$v1);
}

// 2. $obj->{表达式}，例如 $obj->{'b'.'ar'}
$v2 = $obj->{'b'.'ar'};
if ($v2 === 'bar') {
    Log::info('obj->{expr} property ok');
} else {
    Log::fatal('obj->{expr} property fail, got '.$v2);
}

// 3. $obj->{表达式} 作为方法名：$obj->{getName()}()
$methodName = $obj->getName();
$v3 = $obj->{$methodName};
if ($v3 === 'foo') {
    Log::info('obj->{getName()} property ok');
} else {
    Log::fatal('obj->{getName()} property fail, got '.$v3);
}

// 4. 空安全 ?-> {表达式}：$null?->{expr} 应不报错（当前实现为 NullsafeCall 包装）
$null = null;
$tmp = $null?->{'foo'}; // 主要验证解析/运行是否报错
Log::info('null?->{expr} executed');

echo "=== ->{表达式} 测试完成 ===\n";

