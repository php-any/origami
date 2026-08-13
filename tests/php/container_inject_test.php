<?php

namespace tests\php;

/**
 * Container 构造器注入与 Scoped 生命周期测试
 */

interface ContainerInject_LoggerInterface {}

class ContainerInject_FileLogger implements ContainerInject_LoggerInterface {
    public string $channel = 'file';
}

class ContainerInject_UserRepository {
    public function __construct(
        private ContainerInject_LoggerInterface $logger,
    ) {}
    public function channel(): string {
        return $this->logger->channel;
    }
}

class ContainerInject_ScopedService {
    private static int $seq = 0;
    public int $id;
    public function __construct() {
        $this->id = ++self::$seq;
    }
}

$c = new \Container\Container();
$c->bind(ContainerInject_LoggerInterface::class, ContainerInject_FileLogger::class);

$repo = $c->make(ContainerInject_UserRepository::class);
if (!($repo instanceof ContainerInject_UserRepository)) {
    Log::fatal('构造器注入：类型不匹配');
}
if ($repo->channel() !== 'file') {
    Log::fatal('构造器注入：依赖未解析');
}

$c->scoped(ContainerInject_ScopedService::class);
$scope1 = $c->createScope();
$a = $scope1->make(ContainerInject_ScopedService::class);
$b = $scope1->make(ContainerInject_ScopedService::class);
if ($a->id !== $b->id) {
    Log::fatal('Scoped：同一 scope 内应返回同一实例');
}
$scope1->dispose();

$scope2 = $c->createScope();
$c2 = $scope2->make(ContainerInject_ScopedService::class);
if ($c2 === $a) {
    Log::fatal('Scoped：dispose 后新 scope 应创建新实例');
}
$scope2->dispose();

// 根容器不能直接解析 scoped
try {
    $c->make(ContainerInject_ScopedService::class);
    Log::fatal('Scoped：根容器解析应失败');
} catch (\Throwable $e) {
    if (strpos($e->getMessage(), 'scoped service') === false) {
        Log::fatal('Scoped：异常消息不正确: ' . $e->getMessage());
    }
}

Log::info('Container 构造器注入与 Scoped 测试通过');
