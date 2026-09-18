<?php
/**
 * 每次 render before 时打印 listeners['render'] 数量。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

@unlink(storage_path('framework/listener-count.txt'));

Livewire\before('render', function ($target) {
    $bus = app(Livewire\EventBus::class);
    $r = new ReflectionClass($bus);
    $p = $r->getProperty('listeners');
    $p->setAccessible(true);
    $all = $p->getValue($bus);
    $arr = $all['render'] ?? [];
    $n = is_array($arr) ? count($arr) : -1;
    $keys = is_array($arr) ? implode(',', array_map('strval', array_keys($arr))) : '?';
    $cls = is_object($target) ? get_class($target) : gettype($target);
    file_put_contents(
        storage_path('framework/listener-count.txt'),
        "target={$cls} n={$n} keys={$keys}\n",
        FILE_APPEND
    );
});

$http = app(Illuminate\Contracts\Http\Kernel::class);
$req = Illuminate\Http\Request::create('http://127.0.0.1:8000/admin/login', 'GET');
$resp = $http->handle($req);
echo "status=".$resp->getStatusCode()."\n";
echo @file_get_contents(storage_path('framework/listener-count.txt'));
echo "DONE\n";
