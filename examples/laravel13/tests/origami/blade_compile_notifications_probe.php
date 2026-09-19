<?php
/**
 * 用官方 BladeCompiler 编译 database-notifications，检查 raw block 是否还原、use Alignment 是否在结果里。
 */
$_SERVER['HTTP_HOST'] = '127.0.0.1';
$_SERVER['SERVER_NAME'] = '127.0.0.1';
$_SERVER['REQUEST_URI'] = '/';
$_SERVER['REQUEST_METHOD'] = 'GET';
chdir(__DIR__.'/../..');
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$dbg = __DIR__.'/../../storage/origami-debug';
if (!is_dir($dbg)) {
    mkdir($dbg, 0755, true);
}

$path = base_path('vendor/filament/notifications/resources/views/database-notifications.blade.php');
$src = file_get_contents($path);
$compiler = app('blade.compiler');
$out = $compiler->compileString($src);

file_put_contents($dbg.'/blade-compile-notifications.php', $out);

$report = [
    'len' => strlen($out),
    'has_placeholder' => str_contains($out, '@__raw_block_'),
    'has_use_alignment' => str_contains($out, 'use Filament\\Support\\Enums\\Alignment;'),
    'starts' => substr($out, 0, 180),
];
file_put_contents($dbg.'/blade-compile-notifications.json', json_encode($report, JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES | JSON_PRETTY_PRINT));
echo json_encode($report, JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES), "\n";

try {
    $ref = new ReflectionClass(Illuminate\View\AnonymousComponent::class);
    $ctor = $ref->getConstructor();
    $params = (new Illuminate\Support\Collection($ctor->getParameters()))->map->getName()->all();
    echo 'ctor_params='.json_encode($params)."\n";
    $comp = Illuminate\View\AnonymousComponent::resolve([
        'view' => 'filament-panels::components.page.index',
        'data' => [],
    ]);
    echo 'resolved_class='.get_class($comp)."\n";
    echo 'resolved_view='.var_export($comp->render(), true)."\n";
    echo 'resolveView='.var_export($comp->resolveView(), true)."\n";
} catch (Throwable $e) {
    echo 'resolve_ex='.$e->getMessage()."\n".$e->getFile().':'.$e->getLine()."\n";
}
