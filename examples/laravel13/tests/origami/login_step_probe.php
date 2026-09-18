<?php
/**
 * Step-debug Login view render: compile, shouldRender, e(content), engine get.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$panel = Filament\Facades\Filament::getPanel('admin');
Filament\Facades\Filament::setCurrentPanel($panel);
$login = app(Filament\Auth\Pages\Login::class);
$login->mount();

$view = $login->render();
$path = $view->getPath();
$compiler = app('blade.compiler');
$compiler->compile($path);
$compiled = $compiler->getCompiledPath($path);
echo "compiled_exists=".(file_exists($compiled)?'yes':'no')."\n";
echo "compiled_head=".substr(file_get_contents($compiled),0,200)."\n";

$extend = app(Livewire\Mechanisms\ExtendBlade\ExtendBlade::class);
$extend->startLivewireRendering($login);
try {
    // 1) e(content) alone
    ob_start();
    echo e($login->content);
    $slot = ob_get_clean();
    echo "slot_len=".strlen($slot)."\n";
    echo "slot_head=".substr(preg_replace('/\s+/',' ',$slot),0,150)."\n";

    // 2) resolve component
    $component = Illuminate\View\AnonymousComponent::resolve([
        'view' => 'filament-panels::components.page.simple',
        'data' => [],
    ]);
    echo "shouldRender=".($component->shouldRender()?'yes':'no')."\n";
    echo "resolveView=".json_encode($component->resolveView())."\n";

    // 3) engine get
    $engine = app('view.engine.resolver')->resolve('blade');
    echo "engine=".get_class($engine)."\n";
    $__env = app('view');
    $data = array_merge($view->getData(), compact('__env'));
    try {
        $html = $engine->get($path, $data);
        echo "engine_len=".strlen($html)."\n";
        echo "engine_head=".substr(preg_replace('/\s+/',' ',$html),0,300)."\n";
    } catch (Throwable $e) {
        echo "engine_EX=".get_class($e).": ".$e->getMessage()."\n";
        echo "at ".$e->getFile().":".$e->getLine()."\n";
        $p = $e->getPrevious();
        $i=0;
        while ($p && $i<6) {
            echo "prev$i=".get_class($p).": ".$p->getMessage()." @ ".$p->getFile().":".$p->getLine()."\n";
            $p=$p->getPrevious(); $i++;
        }
    }

    // 4) direct bind include
    ob_start();
    try {
        \Closure::bind(function () use ($compiled, $data) {
            extract($data, EXTR_SKIP);
            include $compiled;
        }, $login, $login)();
        $html2 = ob_get_clean();
        echo "include_len=".strlen($html2)."\n";
        echo "include_head=".substr(preg_replace('/\s+/',' ',$html2),0,300)."\n";
    } catch (Throwable $e) {
        $buf = (string)ob_get_clean();
        echo "include_buf=".strlen($buf)."\n";
        echo "include_EX=".get_class($e).": ".$e->getMessage()."\n";
        echo "at ".$e->getFile().":".$e->getLine()."\n";
        $p = $e->getPrevious();
        $i=0;
        while ($p && $i<6) {
            echo "prev$i=".get_class($p).": ".$p->getMessage()." @ ".$p->getFile().":".$p->getLine()."\n";
            $p=$p->getPrevious(); $i++;
        }
    }
} finally {
    $extend->endLivewireRendering();
}
echo "DONE\n";
