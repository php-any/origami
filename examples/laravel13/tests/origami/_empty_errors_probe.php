<?php

require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->bootstrapWith([
    Illuminate\Foundation\Bootstrap\LoadEnvironmentVariables::class,
    Illuminate\Foundation\Bootstrap\LoadConfiguration::class,
    Illuminate\Foundation\Bootstrap\RegisterFacades::class,
    Illuminate\Foundation\Bootstrap\SetRequestForConsole::class,
    Illuminate\Foundation\Bootstrap\RegisterProviders::class,
    Illuminate\Foundation\Bootstrap\BootProviders::class,
]);

function p($m) { file_put_contents(__DIR__.'/_login_progress.txt', $m."\n", FILE_APPEND); }
@unlink(__DIR__.'/_login_progress.txt');

$errors = [];
$v = collect($errors)
    ->filter(function ($value, $key) {
        return false;
    })
    ->toArray();
p('type='.gettype($v));
p('is_array='.var_export(is_array($v), true));
p('count='.count($v));
p('json='.json_encode(['errors' => $v]));
p('json_plain='.json_encode(['errors' => []]));

// Simulate MessageBag empty toArray
$bag = new \Illuminate\Support\MessageBag();
p('bag_json='.json_encode(['errors' => collect($bag->toArray())->filter(fn($v,$k)=>true)->toArray()]));
p('done');
