<?php
/**
 * Check whether attributes survive into an included view file scope.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$dir = storage_path('framework/views');
$file = $dir.'/origami_scope_probe.php';
file_put_contents($file, '<?php
echo "inc_isset=".(isset($attributes)?"yes":"no")."\n";
echo "inc_type=".gettype($attributes ?? null)."\n";
$gdv = get_defined_vars();
echo "inc_gdv_has=".(is_array($gdv) ? (array_key_exists("attributes",$gdv)?"yes":"no") : (isset($gdv["attributes"])?"yes":"no"))."\n";
if (is_object($gdv)) {
  echo "inc_gdv_class=".get_class($gdv)."\n";
  $keys = [];
  foreach ($gdv as $k=>$v) { $keys[] = $k; }
  echo "inc_gdv_keys=".implode(",", $keys)."\n";
} elseif (is_array($gdv)) {
  echo "inc_gdv_keys=".implode(",", array_keys($gdv))."\n";
}
');

$attributes = new Illuminate\View\ComponentAttributeBag(['class'=>'probe']);
$brandName = 'x';
extract(['attributes' => $attributes, 'brandName' => $brandName], EXTR_SKIP);
echo "before_inc_isset=".(isset($attributes)?'yes':'no')."\n";
include $file;
echo "DONE\n";
@unlink($file);
