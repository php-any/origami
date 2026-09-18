<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$bus = app(Livewire\EventBus::class);
$ref = new ReflectionClass($bus);
$prop = $ref->getProperty('listeners');
$prop->setAccessible(true);
$listeners = $prop->getValue($bus);
echo "render_listeners=".count($listeners['render'] ?? [])."\n";

// Wrap startLivewireRendering
$ext = null;
foreach (app()->getBindings() as $k => $_) {
    // skip
}
// Invade ExtendBlade via Livewire mechanisms
$mechanisms = invade(app('livewire'))->mechanisms ?? null;

// Direct static spy
$origComponents = null;
Livewire\Mechanisms\ExtendBlade\ExtendBlade::class;

// Patch via on('render') FIRST priority using before
Livewire\on('render', function ($target) {
    // This registers as normal on - runs in order
});

Livewire\before('render', function ($target, $view) {
    file_put_contents(storage_path('framework/before-render.txt'),
        'before target='.get_class($target)."\n", FILE_APPEND);
});

Livewire\after('render', function ($target, $view) {
    $cur = Livewire\Mechanisms\ExtendBlade\ExtendBlade::currentRendering();
    file_put_contents(storage_path('framework/after-render.txt'),
        'after target='.get_class($target).' current='.(is_object($cur)?get_class($cur):'null')."\n", FILE_APPEND);
});

@unlink(storage_path('framework/before-render.txt'));
@unlink(storage_path('framework/after-render.txt'));
@unlink(storage_path('framework/this-dbg.txt'));

$http = app(Illuminate\Contracts\Http\Kernel::class);
$req = Illuminate\Http\Request::create('http://127.0.0.1:8000/admin/login', 'GET');
$http->handle($req);

echo "---before---\n".@file_get_contents(storage_path('framework/before-render.txt'));
echo "---after---\n".@file_get_contents(storage_path('framework/after-render.txt'));
echo "---this---\n".@file_get_contents(storage_path('framework/this-dbg.txt'));
echo "DONE\n";
