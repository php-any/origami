<?php
/**
 * Dump expression passed into @capture directive during compile.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$seen = [];
Illuminate\Support\Facades\Blade::directive('capture', function (string $expression) use (&$seen): string {
    $seen[] = $expression;
    [$name, $arguments] = str_contains($expression, ',') ?
        array_map('trim', explode(',', $expression, 2)) :
        [$expression, ''];
    $seen[] = 'name='.$name;
    $seen[] = 'args='.$arguments;
    return "
                <?php {$name} = (function (\$args) {
                    return function ({$arguments}) use (\$args) {
                        extract(\$args, EXTR_SKIP);
                        ob_start(); ?>
            ";
});

$blade = <<<'BLADE'
@capture($content, $logo, $isDarkMode = false)
    <div>{{ $logo }}</div>
@endcapture
BLADE;

$compiled = app('blade.compiler')->compileString($blade);
echo "SEEN=".json_encode($seen)."\n";
echo "COMPILED=".json_encode($compiled)."\n";
