<?php

namespace tests\php;

/**
 * Container 核心 API 测试：bind / singleton / instance / make / has / isShared / ServiceProvider
 */

class ContainerCore_ConcreteService {
    public string $tag = 'concrete';
}

class ContainerCore_SingletonService {
    private static int $seq = 0;
    public int $id;

    public function __construct() {
        $this->id = ++self::$seq;
    }
}

interface ContainerCore_RepositoryInterface {}

class ContainerCore_SqlRepository implements ContainerCore_RepositoryInterface {
    public string $driver = 'sql';
}

class ContainerCore_TestProvider extends \Container\ServiceProvider {
    public function register(): void {
        $this->container->singleton(
            ContainerCore_RepositoryInterface::class,
            ContainerCore_SqlRepository::class
        );
    }
}

$c = \Container\Container::getInstance();

// bind + make
$c->bind('demo.service', ContainerCore_ConcreteService::class);
$made = $c->make('demo.service');
if (!($made instanceof ContainerCore_ConcreteService)) {
    Log::fatal('bind/make 失败：类型不匹配');
}
if ($made->tag !== 'concrete') {
    Log::fatal('bind/make 失败：属性不正确');
}

// transient：每次 make 新建
$c->bind(ContainerCore_ConcreteService::class);
$t1 = $c->make(ContainerCore_ConcreteService::class);
$t2 = $c->make(ContainerCore_ConcreteService::class);
if ($t1 === $t2) {
    Log::fatal('transient bind 不应返回同一实例');
}

// singleton 共享实例
$c->singleton(ContainerCore_SingletonService::class);
$a = $c->make(ContainerCore_SingletonService::class);
$b = $c->make(ContainerCore_SingletonService::class);
if ($a->id !== $b->id) {
    Log::fatal('singleton 未返回同一实例');
}
if (!$c->isShared(ContainerCore_SingletonService::class)) {
    Log::fatal('isShared 应对 singleton 返回 true');
}

// instance 预注册
$pre = new ContainerCore_ConcreteService();
$pre->tag = 'prebuilt';
$c->instance('prebuilt.service', $pre);
$got = $c->make('prebuilt.service');
if ($got !== $pre) {
    Log::fatal('instance 未返回预注册对象');
}

// has
if (!$c->has('demo.service')) {
    Log::fatal('has 应对已绑定 abstract 返回 true');
}
if ($c->has('not.registered.at.all')) {
    Log::fatal('has 应对未注册 abstract 返回 false');
}

// ServiceProvider
$local = new \Container\Container();
$local->registerProviders([ContainerCore_TestProvider::class]);
$repo = $local->make(ContainerCore_RepositoryInterface::class);
if (!($repo instanceof ContainerCore_SqlRepository)) {
    Log::fatal('ServiceProvider register 未生效');
}

Log::info('Container 核心 API 测试通过');
