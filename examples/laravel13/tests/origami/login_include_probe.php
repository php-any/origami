<?php
/**
 * Direct Closure::bind include of compiled Login view; surface any throw.
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
$compiled = app('blade.compiler')->getCompiledPath($path);
echo "compiled=$compiled\n";
echo "exists=".(file_exists($compiled)?'yes':'no')."\n";

$extend = app(Livewire\Mechanisms\ExtendBlade\ExtendBlade::class);
$extend->startLivewireRendering($login);
try {
    $__env = app('view');
    $data = array_merge($view->getData(), compact('__env'));
    ob_start();
    try {
        \Closure::bind(function () use ($compiled, $data) {
            extract($data, EXTR_SKIP);
            include $compiled;
        }, $login, $login)();
        $html = ob_get_clean();
        echo "len=".strlen($html)."\n";
        echo "head=".substr(preg_replace('/\s+/',' ', $html),0,400)."\n";
    } catch (Throwable $e) {
        $buf = ob_get_clean();
        echo "buf_len=".strlen((string)$buf)."\n";
        echo "EX=".get_class($e).": ".$e->getMessage()."\n";
        echo "at ".$e->getFile().":".$e->getLine()."\n";
        $p = $e->getPrevious();
        $i = 0;
        while ($p && $i < 5) {
            echo "prev$i=".get_class($p).": ".$p->getMessage()." @ ".$p->getFile().":".$p->getLine()."\n";
            $p = $p->getPrevious();
            $i++;
        }
    }
} finally {
    $extend->endLivewireRendering();
}
echo "DONE\n";
