<?php
/**
 * 检查 render 监听器数量，并在 Notifications 时直接调用 ExtendBlade 同款 push。
 * 同时用 after 确认 ExtendBlade 监听器是否被调用（通过 renderCounter 反射失败则用副作用）。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

@unlink(storage_path('framework/listener-dbg.txt'));

$bus = app(Livewire\EventBus::class);
$r = new ReflectionClass($bus);
foreach (['listeners', 'listenersBefore', 'listenersAfter'] as $name) {
    try {
        $p = $r->getProperty($name);
        $p->setAccessible(true);
        $all = $p->getValue($bus);
        $n = isset($all['render']) ? count($all['render']) : 0;
        file_put_contents(storage_path('framework/listener-dbg.txt'), "{$name}.render count={$n}\n", FILE_APPEND);
    } catch (Throwable $e) {
        file_put_contents(storage_path('framework/listener-dbg.txt'), "{$name} EX: ".$e->getMessage()."\n", FILE_APPEND);
    }
}

// Wrap：在所有已有 on(render) 外包一层日志
try {
    $p = $r->getProperty('listeners');
    $p->setAccessible(true);
    $all = $p->getValue($bus);
    $wrapped = [];
    foreach ($all['render'] ?? [] as $i => $cb) {
        $wrapped[] = function (...$args) use ($cb, $i) {
            $target = $args[0] ?? null;
            $cls = is_object($target) ? get_class($target) : gettype($target);
            $before = Livewire\Mechanisms\ExtendBlade\ExtendBlade::currentRendering();
            $beforeCls = is_object($before) ? get_class($before) : var_export($before, true);
            file_put_contents(storage_path('framework/listener-dbg.txt'), "CALL#{$i} begin target={$cls} current={$beforeCls}\n", FILE_APPEND);
            try {
                $ret = $cb(...$args);
            } catch (Throwable $e) {
                file_put_contents(storage_path('framework/listener-dbg.txt'), "CALL#{$i} EX: ".$e->getMessage()."\n", FILE_APPEND);
                throw $e;
            }
            $after = Livewire\Mechanisms\ExtendBlade\ExtendBlade::currentRendering();
            $afterCls = is_object($after) ? get_class($after) : var_export($after, true);
            file_put_contents(storage_path('framework/listener-dbg.txt'), "CALL#{$i} end current={$afterCls} ret=".gettype($ret)."\n", FILE_APPEND);
            return $ret;
        };
    }
    $all['render'] = $wrapped;
    $p->setValue($bus, $all);
    file_put_contents(storage_path('framework/listener-dbg.txt'), "wrapped=".count($wrapped)."\n", FILE_APPEND);
} catch (Throwable $e) {
    file_put_contents(storage_path('framework/listener-dbg.txt'), "wrap EX: ".$e->getMessage()."\n", FILE_APPEND);
}

$http = app(Illuminate\Contracts\Http\Kernel::class);
$req = Illuminate\Http\Request::create('http://127.0.0.1:8000/admin/login', 'GET');
$resp = $http->handle($req);
echo "status=".$resp->getStatusCode()."\n";
echo "---dbg---\n";
echo @file_get_contents(storage_path('framework/listener-dbg.txt'));
echo "DONE\n";
