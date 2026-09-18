<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

// Manually ensure stack push (simulates what ExtendBlade should do)
Livewire\before('render', function ($target, $view) {
    // Use reflection to call start on the ExtendBlade instance if needed
    Livewire\Mechanisms\ExtendBlade\ExtendBlade::class;
});

// Directly patch: prepend a listener that pushes using the same static API
$ref = new ReflectionClass(Livewire\Mechanisms\ExtendBlade\ExtendBlade::class);
// Call static by invoking via a dummy - use invade on mechanism
$eb = null;
foreach (app(Livewire\LivewireManager::class) /* skip */ ? [] : [] as $x) {}

// Use on with early registration via listeners array unshift
$bus = app(Livewire\EventBus::class);
$r = new ReflectionClass($bus);
$p = $r->getProperty('listeners');
$p->setAccessible(true);
$all = $p->getValue($bus);
$push = function ($target, $view) {
    // Duplicate ExtendBlade push via reflection on static prop
    $ref = new ReflectionClass(Livewire\Mechanisms\ExtendBlade\ExtendBlade::class);
    $prop = $ref->getProperty('livewireComponents');
    $prop->setAccessible(true);
    $stack = $prop->getValue() ?? [];
    if (!is_array($stack)) {
        $stack = [];
    }
    $stack[] = $target;
    $prop->setValue(null, $stack);
    file_put_contents(storage_path('framework/manual-push.txt'),
        'pushed='.get_class($target).' count='.count($stack)."\n", FILE_APPEND);
};
$all['render'] = array_merge([$push], $all['render'] ?? []);
$p->setValue($bus, $all);

@unlink(storage_path('framework/manual-push.txt'));
@unlink(storage_path('framework/this-dbg.txt'));

$http = app(Illuminate\Contracts\Http\Kernel::class);
$req = Illuminate\Http\Request::create('http://127.0.0.1:8000/admin/login', 'GET');
$resp = $http->handle($req);
$c = (string)$resp->getContent();
echo "status=".$resp->getStatusCode()."\n";
echo "has_err=".(str_contains($c,'getBroadcastChannel')?'yes':'no')."\n";
echo "---push---\n".@file_get_contents(storage_path('framework/manual-push.txt'));
echo "---this---\n".@file_get_contents(storage_path('framework/this-dbg.txt'));
echo "DONE\n";
