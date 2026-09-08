<?php

require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->bootstrapWith([
    Illuminate\Foundation\Bootstrap\LoadEnvironmentVariables::class,
    Illuminate\Foundation\Bootstrap\LoadConfiguration::class,
    Illuminate\Foundation\Bootstrap\HandleExceptions::class,
    Illuminate\Foundation\Bootstrap\RegisterFacades::class,
    Illuminate\Foundation\Bootstrap\SetRequestForConsole::class,
    Illuminate\Foundation\Bootstrap\RegisterProviders::class,
    Illuminate\Foundation\Bootstrap\BootProviders::class,
]);

$c = \Illuminate\Support\Carbon::parse('2026-09-07 12:00:00');

$v = 120;
echo "is_numeric=".var_export(is_numeric($v), true)." float=".var_export((float)$v, true)." notfloat=".var_export(!(float)$v, true)."\n";
$neg = -120;
echo "neg is_numeric=".var_export(is_numeric($neg), true)." float=".var_export((float)$neg, true)." notfloat=".var_export(!(float)$neg, true)."\n";

$a1 = $c->copy()->addUnit('minute', -120);
echo "addUnit_minute_neg=". $a1->format('Y-m-d H:i:s')."\n";
$a2 = $c->copy()->addUnit('minute', 120);
echo "addUnit_minute_pos=". $a2->format('Y-m-d H:i:s')."\n";
$a3 = $c->copy()->subUnit('minute', 120);
echo "subUnit_minute=". $a3->format('Y-m-d H:i:s')."\n";
$a4 = $c->copy()->sub('minute', 120);
echo "sub_minute_120=". $a4->format('Y-m-d H:i:s')."\n";
$a5 = $c->copy()->addMinutes(120);
echo "addMinutes=". $a5->format('Y-m-d H:i:s')."\n";
$a6 = $c->copy()->subMinutes(120);
echo "subMinutes=". $a6->format('Y-m-d H:i:s')."\n";

// What does __call resolve?
$ref = new ReflectionClass($c);
echo "has_subMinutes_method=". ($ref->hasMethod('subMinutes') ? 'yes' : 'no')."\n";
