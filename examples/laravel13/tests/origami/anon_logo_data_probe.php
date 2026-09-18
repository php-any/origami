<?php
/**
 * Inspect AnonymousComponent data() for logo.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$factory = app('view');
$resolver = app('view.finder');
// Resolve component class via Blade
$comp = new Illuminate\View\AnonymousComponent(
    'filament-panels::components.logo',
    []
);
$data = $comp->data();
echo "data_type=".gettype($data)."\n";
if (is_array($data)) {
    echo "keys=".implode(',', array_keys($data))."\n";
    echo "has_attributes=".(array_key_exists('attributes', $data)?'yes':'no')."\n";
    if (isset($data['attributes'])) {
        echo "attr_class=".get_class($data['attributes'])."\n";
    } else {
        echo "attr_missing_or_null\n";
        echo "attr_isset=".(isset($data['attributes'])?'yes':'no')."\n";
        var_export(array_key_exists('attributes', $data));
        echo "\n";
    }
} elseif (is_object($data)) {
    echo "data_class=".get_class($data)."\n";
    echo "aa_attr=".(isset($data['attributes'])?'yes':'no')."\n";
}

// Simulate view extract path
extract($data, EXTR_SKIP);
echo "after_extract_isset=".(isset($attributes)?'yes':'no')."\n";
echo "after_type=".gettype($attributes ?? null)."\n";

$gdv = get_defined_vars();
echo "gdv_has=".(isset($gdv['attributes']) || (is_array($gdv) && array_key_exists('attributes', $gdv)) || (is_object($gdv) && isset($gdv['attributes'])) ? 'yes':'no')."\n";

echo "DONE\n";
