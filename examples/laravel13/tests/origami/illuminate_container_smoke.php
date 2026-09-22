<?php

namespace tests\origami;

/**
 * Container spike：ORIGAMI_STD_CONTAINER=1 时 Go Container 可 bind/make；
 * Application extends Container 在默认关闭下仍走 vendor PHP。
 */

use Illuminate\Container\Container;

$c = new Container();
$c->instance('app.name', 'Origami');
if ($c->make('app.name') !== 'Origami') {
    \Log::fatal('instance/make 失败');
}

$c->singleton('counter', function () {
    static $n = 0;
    $n++;
    return $n;
});
$a = $c->make('counter');
$b = $c->make('counter');
if ($a !== $b) {
    \Log::fatal('singleton 应共享实例');
}

if (!$c->bound('app.name')) {
    \Log::fatal('bound 失败');
}

// spike：Application 是否能继承（仅当 Go Container 启用且 Application 尚未加载）
$extendsOk = is_subclass_of(\Illuminate\Foundation\Application::class, Container::class);
\Log::info('Application extends Container: ' . ($extendsOk ? 'yes' : 'no'));
\Log::info('illuminate_container_spike 测试通过');
