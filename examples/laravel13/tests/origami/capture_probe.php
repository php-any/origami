<?php
/**
 * Reproduce Filament @capture directive compilation under Origami.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$expression = '$content, $logo, $isDarkMode = false';
[$name, $arguments] = str_contains($expression, ',') ?
    array_map('trim', explode(',', $expression, 2)) :
    [$expression, ''];

echo "name=".json_encode($name)."\n";
echo "arguments=".json_encode($arguments)."\n";

$out = "
                <?php {$name} = (function (\$args) {
                    return function ({$arguments}) use (\$args) {
                        extract(\$args, EXTR_SKIP);
                        ob_start(); ?>
            ";
echo "out=".json_encode($out)."\n";
echo (str_contains($out, '$content =') ? "HAS_CONTENT\n" : "MISSING_CONTENT\n");

// Full blade compile
$blade = <<<'BLADE'
@capture($content, $logo, $isDarkMode = false)
    <div>{{ $logo }}</div>
@endcapture
{{ $content('x') }}
BLADE;

$compiled = app('blade.compiler')->compileString($blade);
echo "---COMPILED---\n";
echo $compiled, "\n";
echo "---\n";
echo (str_contains($compiled, '$content =') ? "COMPILED_OK\n" : "COMPILED_BAD\n");
echo (preg_match('/<\?php\s+=/', $compiled) ? "EMPTY_ASSIGN\n" : "no_empty_assign\n");
