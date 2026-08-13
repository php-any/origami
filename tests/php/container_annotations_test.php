<?php

namespace tests\php;

/**
 * Container 注解与 resolver 增强测试：Component / Singleton / Bind
 */

$c = new \Container\Container();
$c->scan(__DIR__ . '/container_fixtures/phase2');

// 类型提示构造器注入 + #[Bind]
$order = $c->make(\tests\php\container_fixtures\Phase2Anno_OrderService::class);
if ($order->mailerTag() !== 'mailer') {
    Log::fatal('构造器注入：mailer 未解析');
}
if ($order->loggerChannel() !== 'file') {
    Log::fatal('Bind + 构造器注入：logger 未解析');
}

// #[Singleton]
$s1 = $c->make(\tests\php\container_fixtures\Phase2Anno_Config::class);
$s2 = $c->make(\tests\php\container_fixtures\Phase2Anno_Config::class);
if ($s1->id !== $s2->id) {
    Log::fatal('Singleton 注解：应返回同一实例');
}

// #[Component] transient（按别名注册）
$m1 = $c->make('phase2.transient');
$m2 = $c->make('phase2.transient');
if ($m1 === $m2) {
    Log::fatal('Component 注解：应为 transient');
}

// 参数名绑定（无类型提示时按参数名解析）
$c->bind('loggerApp', \tests\php\container_fixtures\Phase2Anno_FileLogger::class);
class ContainerAnno_LoggerByName {
    private $loggerApp;
    public function __construct($loggerApp) {
        $this->loggerApp = $loggerApp;
    }
    public function channel(): string {
        return $this->loggerApp->channel;
    }
}
$byName = $c->make(ContainerAnno_LoggerByName::class);
if ($byName->channel() !== 'file') {
    Log::fatal('参数名绑定：loggerApp 未解析');
}

// 可选依赖：有默认值时容器无 binding 不报错
class ContainerAnno_OptionalDep {
    public string $tag;
    public function __construct(string $tag = 'default') {
        $this->tag = $tag;
    }
}
$opt = $c->make(ContainerAnno_OptionalDep::class);
if ($opt->tag !== 'default') {
    Log::fatal('可选依赖：应使用默认值');
}

Log::info('Container 注解与 resolver 测试通过');
