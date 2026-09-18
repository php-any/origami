<?php
/**
 * Render filament logo view and dump whether $attributes is in get_defined_vars / after extract.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$view = view('filament-panels::components.logo', [
    'attributes' => new Illuminate\View\ComponentAttributeBag(['class' => 'probe']),
]);

try {
    $html = $view->render();
    echo "render_ok len=".strlen($html)."\n";
    echo "snippet=".substr(preg_replace('/\s+/', ' ', $html), 0, 200)."\n";
} catch (Throwable $e) {
    echo "EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
}

// Mirror capture extract path like compiled logo view
$attributes = new Illuminate\View\ComponentAttributeBag(['class' => 'probe']);
$brandName = 'Brand';
$getLogoClasses = fn (bool $isDarkMode): string => 'fi-logo';
$logoStyles = 'height: 1.5rem';

$content = (function ($args) {
    return function ($logo, $isDarkMode = false) use ($args) {
        $keys = is_array($args) ? array_keys($args) : (is_object($args) ? array_keys((array)$args) : []);
        echo "args_has_attributes=".(array_key_exists('attributes', (array)$args) || (is_array($args) && array_key_exists('attributes', $args)) ? 'yes' : 'no')."\n";
        echo "args_type=".gettype($args)."\n";
        if (is_array($args)) {
            echo "args_keys=".implode(',', array_keys($args))."\n";
        } elseif (is_object($args)) {
            echo "args_class=".get_class($args)."\n";
            // ArrayAccess?
            if ($args instanceof ArrayAccess) {
                echo "args_aa_attributes=".(isset($args['attributes']) ? gettype($args['attributes']) : 'missing')."\n";
            }
            $vars = get_object_vars($args);
            echo "args_obj_keys=".implode(',', array_keys($vars))."\n";
        }
        extract($args, EXTR_SKIP);
        echo "after_extract_isset=".(isset($attributes) ? 'yes' : 'no')."\n";
        echo "after_extract_type=".gettype($attributes ?? null)."\n";
        if (isset($attributes) && is_object($attributes)) {
            echo "after_class=".get_class($attributes)."\n";
            try {
                echo "html=".$attributes->class(['x'])."\n";
            } catch (Throwable $e) {
                echo "class_EX: ".$e->getMessage()."\n";
            }
        }
        return 'ok';
    };
})(get_defined_vars());

echo $content(null);
echo "DONE\n";
