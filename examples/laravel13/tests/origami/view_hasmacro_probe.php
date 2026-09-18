<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$kernel = $app->make(Illuminate\Contracts\Console\Kernel::class);
$kernel->bootstrap();
echo "boot\n";
echo "method_exists hasMacro: ".(method_exists(Illuminate\View\View::class,'hasMacro')?'yes':'no')."\n";
try {
  $r = Illuminate\View\View::hasMacro('x');
  echo "hasMacro result: ";
  var_dump($r);
} catch (Throwable $e) {
  echo "EX: ".$e->getMessage()."\n";
}
// Reflect whether static::$macros exists via a closure bound to View
$fn = function() {
  return array_key_exists('macros', get_class_vars(Illuminate\View\View::class));
};
echo "has macros prop: ";
var_dump($fn());
