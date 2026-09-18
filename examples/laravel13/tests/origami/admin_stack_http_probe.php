<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$kernel = $app->make(Illuminate\Contracts\Console\Kernel::class);
$kernel->bootstrap();

@unlink(storage_path('framework/stack-probe.txt'));

Livewire\before('render', function ($target, $view) {
    $cls = is_object($target) ? get_class($target) : gettype($target);
    $cur = Livewire\Mechanisms\ExtendBlade\ExtendBlade::currentRendering();
    $curCls = is_object($cur) ? get_class($cur) : var_export($cur, true);
    $viewName = is_object($view) && method_exists($view, 'getName') ? $view->getName() : '?';
    file_put_contents(
        storage_path('framework/stack-probe.txt'),
        "BEFORE target={$cls} current={$curCls} view={$viewName}\n",
        FILE_APPEND
    );
});

Livewire\on('render', function ($target, $view) {
    $cls = is_object($target) ? get_class($target) : gettype($target);
    $cur = Livewire\Mechanisms\ExtendBlade\ExtendBlade::currentRendering();
    $curCls = is_object($cur) ? get_class($cur) : var_export($cur, true);
    $viewName = is_object($view) && method_exists($view, 'getName') ? $view->getName() : '?';
    file_put_contents(
        storage_path('framework/stack-probe.txt'),
        "ON target={$cls} current={$curCls} view={$viewName}\n",
        FILE_APPEND
    );
    return function ($html) use ($cls) {
        $cur = Livewire\Mechanisms\ExtendBlade\ExtendBlade::currentRendering();
        $curCls = is_object($cur) ? get_class($cur) : var_export($cur, true);
        file_put_contents(
            storage_path('framework/stack-probe.txt'),
            "FINISH targetWas={$cls} current={$curCls}\n",
            FILE_APPEND
        );
        return $html;
    };
});

try {
    $http = $app->make(Illuminate\Contracts\Http\Kernel::class);
    $req = Illuminate\Http\Request::create('http://127.0.0.1:8000/admin/login', 'GET');
    echo "handling\n";
    $response = $http->handle($req);
    echo "status=".$response->getStatusCode()."\n";
    $c = (string)$response->getContent();
    echo "len=".strlen($c)."\n";
    echo "has_err=".(str_contains($c, 'getBroadcastChannel') ? 'yes' : 'no')."\n";
} catch (Throwable $e) {
    echo "EX: ".$e->getMessage()."\n";
}

echo "---stack---\n";
echo @file_get_contents(storage_path('framework/stack-probe.txt'));
echo "DONE\n";
