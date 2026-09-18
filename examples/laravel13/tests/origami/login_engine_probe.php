<?php
/**
 * Dig into why Livewire login view renders empty.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

try {
    $panel = Filament\Facades\Filament::getPanel('admin');
    Filament\Facades\Filament::setCurrentPanel($panel);
    $login = app(Filament\Auth\Pages\Login::class);
    $login->mount();
    echo "login=".get_class($login)."\n";
    echo "has_content_method=".(method_exists($login,'getContent')?'yes':'no')."\n";

    // content property via __get
    try {
        $c = $login->content;
        echo "content_type=".gettype($c).(is_object($c)?' '.get_class($c):'')."\n";
    } catch (Throwable $e) {
        echo "content_err=".$e->getMessage()."\n";
    }

    $extend = app(Livewire\Mechanisms\ExtendBlade\ExtendBlade::class);
    $extend->startLivewireRendering($login);
    try {
        $view = $login->render();
        $path = $view->getPath();
        echo "view_path=$path\n";
        echo "is_lw=". (Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()?'yes':'no')."\n";

        // Direct evaluate via engine like Livewire does
        $engine = app('view.engine.resolver')->resolve('blade');
        echo "engine=".get_class($engine)."\n";
        $data = array_merge($view->getData(), ['__env' => app('view')]);
        try {
            $html = $engine->get($path, $data);
            echo "engine_len=".strlen($html)."\n";
            echo "engine_head=".substr(preg_replace('/\s+/',' ', $html),0,300)."\n";
        } catch (Throwable $e) {
            echo "engine_EX=".get_class($e).": ".$e->getMessage()."\n";
            echo "at ".$e->getFile().":".$e->getLine()."\n";
            $p = $e->getPrevious();
            while ($p) {
                echo "prev=".get_class($p).": ".$p->getMessage()." @ ".$p->getFile().":".$p->getLine()."\n";
                $p = $p->getPrevious();
            }
        }
    } finally {
        $extend->endLivewireRendering();
    }
} catch (Throwable $e) {
    echo "EX: ".get_class($e).": ".$e->getMessage()."\n";
    echo "at ".$e->getFile().":".$e->getLine()."\n";
}
echo "DONE\n";
