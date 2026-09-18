<?php

echo "PROBE start\n";

use Illuminate\Foundation\Bootstrap\BootProviders;
use Illuminate\Foundation\Bootstrap\HandleExceptions;
use Illuminate\Foundation\Bootstrap\LoadConfiguration;
use Illuminate\Foundation\Bootstrap\LoadEnvironmentVariables;
use Illuminate\Foundation\Bootstrap\RegisterFacades;
use Illuminate\Foundation\Bootstrap\RegisterProviders;
use Illuminate\Foundation\Bootstrap\SetRequestForConsole;
use Illuminate\Http\Request;

require __DIR__.'/../../vendor/autoload.php';

function probe_log(string $msg): void
{
    echo 'PROBE '.date('H:i:s').' '.$msg."\n";
}

/** @var Illuminate\Foundation\Application $app */
$app = require __DIR__.'/../../bootstrap/app.php';
$app->bootstrapWith([
    LoadEnvironmentVariables::class,
    LoadConfiguration::class,
    HandleExceptions::class,
    RegisterFacades::class,
    SetRequestForConsole::class,
    RegisterProviders::class,
    BootProviders::class,
]);

$request = Request::create('/admin/login', 'GET');
$app->instance('request', $request);
$router = $app->make('router');
probe_log('match');
$route = $router->getRoutes()->match($request);
probe_log('uri='.$route->uri().' action='.$route->getActionName());
$mw = $router->gatherRouteMiddleware($route);
probe_log('middleware n='.count($mw));
foreach ($mw as $i => $class) {
    probe_log('mw['.$i.']='.$class);
}

probe_log('route run');
$raw = $route->run();
probe_log('ran type='.get_debug_type($raw));

echo "OK: admin_route_probe\n";
