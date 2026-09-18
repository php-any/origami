<?php
/**
 * Probe extract(get_defined_vars) and attributes inside capture-like closure.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$attributes = new Illuminate\View\ComponentAttributeBag(['class' => 'x']);
$brandName = 'Test';
$logo = null;

$content = (function ($args) {
    return function ($logo, $isDarkMode = false) use ($args) {
        extract($args, EXTR_SKIP);
        echo "attrs_isset=".(isset($attributes)?'yes':'no')."\n";
        echo "attrs_type=".gettype($attributes ?? null)."\n";
        if (isset($attributes) && is_object($attributes)) {
            echo "attrs_class=".get_class($attributes)."\n";
            echo "html=".$attributes->class(['y'])."\n";
        }
        return 'ok';
    };
})(get_defined_vars());

echo $content(null);
echo "DONE\n";
