<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

// Exact copy of Filament's directive body in THIS file
$mirror = function (string $expression): string {
    [$name, $arguments] = str_contains($expression, ',') ?
        array_map('trim', explode(',', $expression, 2)) :
        [$expression, ''];

    return "
                <?php {$name} = (function (\$args) {
                    return function ({$arguments}) use (\$args) {
                        extract(\$args, EXTR_SKIP);
                        ob_start(); ?>
            ";
};

$expr = '$content, $logo, $isDarkMode = false';
$mout = $mirror($expr);
echo "mirror_ok=".(str_contains($mout, '$content =')?'yes':'no')."\n";

$compiler = app('blade.compiler');
$ref = new ReflectionProperty($compiler, 'customDirectives');
$ref->setAccessible(true);
$ffn = $ref->getValue($compiler)['capture'];
$fout = $ffn($expr);
echo "filament_ok=".(str_contains($fout, '$content =')?'yes':'no')."\n";
echo "filament_ret=".json_encode($fout)."\n";

// Check whether Filament closure's $name is actually assigned: wrap by parsing source
$path = __DIR__.'/../../vendor/filament/support/src/SupportServiceProvider.php';
$lines = file($path);
for ($i = 172; $i < 185; $i++) {
    echo ($i+1).":".$lines[$i];
}
