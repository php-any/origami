<?php

use Illuminate\Http\Request;

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

$payload = '{"components":[{"snapshot":"{\\"data\\":{}}","updates":{"email":"a"},"calls":[]}]}';
$req = Request::create('/livewire-6dd39ca7/update', 'POST', [], [], [], [
    'CONTENT_TYPE' => 'application/json',
    'HTTP_ACCEPT' => 'application/json',
    'HTTP_X_LIVEWIRE' => 'true',
], $payload);

echo "class=".get_class($req)."\n";
echo "isJson=".var_export($req->isJson(), true)."\n";
echo "getContent_len=".strlen($req->getContent())."\n";
echo "getContent=".substr($req->getContent(), 0, 60)."\n";

// Force json() bag
$json = $req->json();
echo "json_class=".get_class($json)."\n";
echo "json_all=".json_encode($json->all())."\n";
echo "input_components=".json_encode($req->input('components'))."\n";
echo "all=".json_encode($req->all())."\n";
