<?php
/**
 * Mirror logo @capture extract without rendering Filament view.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$attributes = new Illuminate\View\ComponentAttributeBag(['class' => 'probe']);
$brandName = 'Brand';
$getLogoClasses = fn (bool $isDarkMode): string => 'fi-logo';
$logoStyles = 'height: 1.5rem';

$gdv = get_defined_vars();
echo "gdv_type=".gettype($gdv)."\n";
if (is_array($gdv)) {
    echo "gdv_has_attributes=".(array_key_exists('attributes', $gdv)?'yes':'no')."\n";
    echo "gdv_attr_type=".gettype($gdv['attributes'] ?? null)."\n";
} elseif (is_object($gdv)) {
    echo "gdv_class=".get_class($gdv)."\n";
    echo "gdv_aa=".(isset($gdv['attributes'])?'yes':'no')."\n";
}

$content = (function ($args) {
    return function ($logo, $isDarkMode = false) use ($args) {
        echo "inner_args_type=".gettype($args)."\n";
        if (is_array($args)) {
            echo "inner_has=".(array_key_exists('attributes', $args)?'yes':'no')."\n";
        } elseif ($args instanceof ArrayAccess) {
            echo "inner_aa_has=".(isset($args['attributes'])?'yes':'no')."\n";
            if (isset($args['attributes'])) {
                echo "inner_aa_type=".gettype($args['attributes'])."\n";
            }
        }
        extract($args, EXTR_SKIP);
        echo "after_isset=".(isset($attributes)?'yes':'no')."\n";
        echo "after_type=".gettype($attributes ?? null)."\n";
        if (is_object($attributes ?? null)) {
            echo "after_class=".get_class($attributes)."\n";
            echo "html=".$attributes->class(['x'])."\n";
        }
        return 'ok';
    };
})(get_defined_vars());

echo $content(null)."\n";
echo "DONE\n";
