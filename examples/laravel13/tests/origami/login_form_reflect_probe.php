<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$m = new ReflectionMethod(Filament\Auth\Pages\Login::class, 'form');
echo "params: ".$m->getNumberOfParameters()."\n";
$ps = $m->getParameters();
echo "count: ".count($ps)."\n";
if (count($ps)>0) {
  $p = $ps[0];
  echo "name: ".$p->getName()."\n";
  $t = $p->getType();
  echo "type class: ".($t === null ? 'null' : get_class($t))."\n";
  if ($t) {
    echo "type name: ".$t->getName()."\n";
    echo "is named: ".(($t instanceof ReflectionNamedType)?'yes':'no')."\n";
    echo "is_a Schema: ".(is_a($t->getName(), Filament\Schemas\Schema::class, true)?'yes':'no')."\n";
  }
}
