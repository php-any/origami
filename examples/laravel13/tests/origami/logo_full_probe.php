<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

try {
    $html = Illuminate\View\Compilers\BladeCompiler::renderComponent(
        new Illuminate\View\AnonymousComponent('origami-logo-full', [])
    );
    echo "ok len=".strlen($html)."\n";
    echo "OUT:\n".$html."\n";
} catch (Throwable $e) {
    echo "EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
}
echo "DONE\n";
