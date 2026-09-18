<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$m = new ReflectionMethod(Filament\Auth\Pages\Login::class, 'form');
$t = $m->getParameters()[0]->getType();
$type = $t->getName();
$schema = Filament\Schemas\Schema::class;
echo "type=[$type]\n";
echo "schema=[$schema]\n";
echo "eq=".($type === $schema ? 'yes':'no')."\n";
echo "is_a pos=". (is_a($type, $schema, true) ? 'yes':'no') ."\n";
echo "is_a named=". (is_a($type, $schema, allow_string: true) ? 'yes':'no') ."\n";
echo "class_exists type=". (class_exists($type)?'yes':'no') ."\n";
echo "class_exists schema=". (class_exists($schema)?'yes':'no') ."\n";
$obj = new $type(null);
echo "instanceof=". (($obj instanceof Filament\Schemas\Schema)?'yes':'no') ."\n";
