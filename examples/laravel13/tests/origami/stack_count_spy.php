<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

Livewire\on('render', function ($target, $view) {
    // Runs after ExtendBlade in registration order if we register now — actually at end
    $ref = new ReflectionClass(Livewire\Mechanisms\ExtendBlade\ExtendBlade::class);
    $prop = $ref->getProperty('livewireComponents');
    $prop->setAccessible(true);
    $stack = $prop->getValue();
    $count = is_countable($stack) ? count($stack) : -1;
    $cur = Livewire\Mechanisms\ExtendBlade\ExtendBlade::currentRendering();
    file_put_contents(storage_path('framework/stack-spy.txt'),
        'target='.get_class($target).
        ' stackCount='.$count.
        ' current='.(is_object($cur)?get_class($cur):json_encode($cur)).
        ' stackType='.gettype($stack)."\n",
        FILE_APPEND
    );
});

@unlink(storage_path('framework/stack-spy.txt'));

$http = app(Illuminate\Contracts\Http\Kernel::class);
$req = Illuminate\Http\Request::create('http://127.0.0.1:8000/admin/login', 'GET');
$http->handle($req);
echo file_get_contents(storage_path('framework/stack-spy.txt'));
echo "DONE\n";
