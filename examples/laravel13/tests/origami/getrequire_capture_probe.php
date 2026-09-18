<?php
/**
 * Exact Filesystem::getRequire + @capture get_defined_vars mirror.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$dir = storage_path('framework/views');
$file = $dir.'/origami_gdv_probe.php';
file_put_contents($file, <<<'PHP'
<?php
$gdv = get_defined_vars();
echo "gdv_type=".gettype($gdv)."\n";
if (is_array($gdv)) {
    echo "keys=".implode(',', array_keys($gdv))."\n";
    echo "has_attr=".(array_key_exists('attributes', $gdv)?'yes':'no')."\n";
    echo "attr_type=".gettype($gdv['attributes'] ?? null)."\n";
} elseif (is_object($gdv)) {
    echo "gdv_class=".get_class($gdv)."\n";
    echo "has_attr=".(isset($gdv['attributes'])?'yes':'no')."\n";
    echo "attr_type=".gettype($gdv['attributes'] ?? null)."\n";
}
echo "direct_isset=".(isset($attributes)?'yes':'no')."\n";

$content = (function ($args) {
    return function ($logo) use ($args) {
        echo "args_type=".gettype($args)."\n";
        if (is_array($args)) {
            echo "args_has=".(array_key_exists('attributes',$args)?'yes':'no')."\n";
        } elseif (is_object($args)) {
            echo "args_has=".(isset($args['attributes'])?'yes':'no')."\n";
            echo "args_attr_type=".gettype($args['attributes'] ?? null)."\n";
        }
        extract($args, EXTR_SKIP);
        echo "after_isset=".(isset($attributes)?'yes':'no')."\n";
        echo "after_type=".gettype($attributes ?? null)."\n";
        return 'ok';
    };
})(get_defined_vars());
echo $content(null)."\n";
PHP);

$data = [
    'attributes' => new Illuminate\View\ComponentAttributeBag(['class' => 'probe']),
    'brandName' => 'Brand',
];

$__path = $file;
$__data = $data;
(static function () use ($__path, $__data) {
    extract($__data, EXTR_SKIP);
    echo "static_after_extract_isset=".(isset($attributes)?'yes':'no')."\n";
    echo "static_attr_type=".gettype($attributes ?? null)."\n";
    return require $__path;
})();

echo "DONE\n";
@unlink($file);
