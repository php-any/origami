<?php
/**
 * Call BladeCompiler::renderComponent on AnonymousComponent for logo.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$comp = new Illuminate\View\AnonymousComponent('filament-panels::components.logo', []);
$data = $comp->data();
echo "data_keys=".implode(',', array_keys($data))."\n";
echo "data_attr=".(isset($data['attributes'])?get_class($data['attributes']):'no')."\n";

try {
    $html = Illuminate\View\Compilers\BladeCompiler::renderComponent($comp);
    echo "ok len=".strlen($html)."\n";
    echo substr(preg_replace('/\s+/',' ',$html),0,200)."\n";
} catch (Throwable $e) {
    echo "EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
}
echo "DONE\n";
