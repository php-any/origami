<?php
/**
 * LoginResponse instanceof Illuminate\Contracts\Support\Responsable
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$r = app(\Filament\Auth\Http\Responses\LoginResponse::class);
echo "class=".get_class($r)."\n";
echo "impl_filament=".( $r instanceof \Filament\Auth\Http\Responses\Contracts\LoginResponse ? 'yes' : 'no')."\n";
echo "impl_illuminate=".( $r instanceof \Illuminate\Contracts\Support\Responsable ? 'yes' : 'no')."\n";

$ifaces = class_implements($r);
echo "implements=".implode(',', array_keys($ifaces ?: []))."\n";
echo "DONE\n";
