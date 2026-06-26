<?php

namespace tests\php;

/**
 * Container Phase 3：alias / 工厂闭包 / #[Bind] 注解
 */

interface ContainerP3_CacheInterface {
    public function get(string $key): ?string;
}

class ContainerP3_FileCache implements ContainerP3_CacheInterface {
    public function get(string $key): ?string {
        return 'file:' . $key;
    }
}

class ContainerP3_ConfigService {
    public string $source = 'factory';
}

class ContainerP3_Consumer {
    public function __construct(private ContainerP3_CacheInterface $cache) {}
    public function read(string $key): ?string {
        return $this->cache->get($key);
    }
}

$c = new \Container\Container();

// alias
$c->bind(ContainerP3_CacheInterface::class, ContainerP3_FileCache::class);
$c->alias('cache', ContainerP3_CacheInterface::class);
if (!$c->has('cache')) {
    Log::fatal('alias：has 应对 alias 返回 true');
}
$viaAlias = $c->make('cache');
if (!($viaAlias instanceof ContainerP3_FileCache)) {
    Log::fatal('alias：make 未解析到实现');
}
if ($viaAlias->get('k') !== 'file:k') {
    Log::fatal('alias：实例行为不正确');
}

// 工厂闭包 bind（transient）
$c->bind('config.service', function (\Container\Container $container) {
    $svc = new ContainerP3_ConfigService();
    $svc->source = 'closure';
    return $svc;
});
$f1 = $c->make('config.service');
$f2 = $c->make('config.service');
if ($f1->source !== 'closure') {
    Log::fatal('工厂闭包：source 不正确');
}
if ($f1 === $f2) {
    Log::fatal('工厂 bind 应为 transient，每次新建');
}

// 工厂闭包 singleton
$c->singleton('config.singleton', function (\Container\Container $container) {
    return new ContainerP3_ConfigService();
});
$s1 = $c->make('config.singleton');
$s2 = $c->make('config.singleton');
if ($s1 !== $s2) {
    Log::fatal('工厂 singleton 应返回同一实例');
}

// 构造器注入 + alias
$c2 = new \Container\Container();
$c2->bind(ContainerP3_CacheInterface::class, ContainerP3_FileCache::class);
$c2->alias('cache', ContainerP3_CacheInterface::class);
$consumer = $c2->make(ContainerP3_Consumer::class);
if ($consumer->read('x') !== 'file:x') {
    Log::fatal('alias + 构造器注入失败');
}

// #[Bind] 注解（scan 触发类加载时注册）
$c3 = new \Container\Container();
$c3->scan(__DIR__ . '/container_fixtures/phase3');
$annotated = $c3->make(\tests\php\container_fixtures\Phase3Bind_CacheInterface::class);
if (!($annotated instanceof \tests\php\container_fixtures\Phase3Bind_AnnotatedCache)) {
    Log::fatal('#[Bind] 注解未注册实现');
}
if ($annotated->get('y') !== 'annotated:y') {
    Log::fatal('#[Bind] 注解实例行为不正确');
}

Log::info('Container Phase 3 测试通过');
