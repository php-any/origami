<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$bus = app(Livewire\EventBus::class);
$bus->on('render', function ($target, $view) {
    $cls = is_object($target) ? get_class($target) : gettype($target);
    $cur = Livewire\Mechanisms\ExtendBlade\ExtendBlade::currentRendering();
    $curCls = is_object($cur) ? get_class($cur) : 'null';
    file_put_contents(
        storage_path('framework/render-dbg.txt'),
        "render target=$cls currentBeforeFinish=$curCls view=".($view->getName()??'?')."\n",
        FILE_APPEND
    );
});

@unlink(storage_path('framework/render-dbg.txt'));
@unlink(storage_path('framework/this-dbg.txt'));

$http = app(Illuminate\Contracts\Http\Kernel::class);
$req = Illuminate\Http\Request::create('http://127.0.0.1:8000/admin/login', 'GET');
$response = $http->handle($req);
echo "status=".$response->getStatusCode()." len=".strlen((string)$response->getContent())."\n";
echo "---render-dbg---\n";
echo @file_get_contents(storage_path('framework/render-dbg.txt'));
echo "---this-dbg---\n";
echo @file_get_contents(storage_path('framework/this-dbg.txt'));
echo "DONE\n";
