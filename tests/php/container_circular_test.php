<?php

namespace tests\php;

class ContainerCircular_A {
    public function __construct(private \tests\php\ContainerCircular_B $b) {}
}

class ContainerCircular_B {
    public function __construct(private \tests\php\ContainerCircular_A $a) {}
}

$c = new \Container\Container();
$c->bind(\tests\php\ContainerCircular_A::class);
$c->bind(\tests\php\ContainerCircular_B::class);

try {
    $c->make(\tests\php\ContainerCircular_A::class);
    Log::fatal('循环依赖：应抛出异常');
} catch (\Throwable $e) {
    $msg = $e->getMessage();
    if (strpos($msg, 'Circular dependency') === false) {
        Log::fatal('循环依赖：异常消息不正确: ' . $msg);
    }
}

Log::info('Container 循环依赖检测测试通过');
