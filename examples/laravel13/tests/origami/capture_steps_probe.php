<?php
/**
 * Step through Filament @capture body and list assignment.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$expression = '$content, $logo, $isDarkMode = false';
echo "contains_comma=".(str_contains($expression, ',')?'yes':'no')."\n";
$exploded = explode(',', $expression, 2);
echo "exploded=".json_encode($exploded)." gettype=".gettype($exploded)."\n";
$mapped = array_map('trim', $exploded);
echo "mapped=".json_encode($mapped)." gettype=".gettype($mapped)." is_array=".(is_array($mapped)?'yes':'no')."\n";

[$name, $arguments] = $mapped;
echo "name=".json_encode($name)." args=".json_encode($arguments)."\n";

// Exact Filament return template
$out = "
                <?php {$name} = (function (\$args) {
                    return function ({$arguments}) use (\$args) {
                        extract(\$args, EXTR_SKIP);
                        ob_start(); ?>
            ";
echo "out_ok=".(str_contains($out, '$content =')?'yes':'no')."\n";

// Call Filament's stored closure
$compiler = app('blade.compiler');
$ref = new ReflectionProperty($compiler, 'customDirectives');
$ref->setAccessible(true);
$fn = $ref->getValue($compiler)['capture'];
echo "fn_class=".get_class($fn)."\n";
// Reflect closure static variables / use
$rf = new ReflectionFunction($fn);
echo "fn_params=".json_encode(array_map(fn($p)=>$p->getName(), $rf->getParameters()))."\n";
$ret = $fn($expression);
echo "fn_ret=".json_encode($ret)."\n";
echo "fn_ok=".(str_contains($ret, '$content =')?'yes':'no')."\n";

// Inline recreate identical closure
$clone = function (string $expression): string {
    [$name, $arguments] = str_contains($expression, ',') ?
        array_map('trim', explode(',', $expression, 2)) :
        [$expression, ''];
    echo "clone_name=".json_encode($name)." clone_args=".json_encode($arguments)."\n";
    return "
                <?php {$name} = (function (\$args) {
                    return function ({$arguments}) use (\$args) {
                        extract(\$args, EXTR_SKIP);
                        ob_start(); ?>
            ";
};
$cret = $clone($expression);
echo "clone_ok=".(str_contains($cret, '$content =')?'yes':'no')."\n";
