<?php

use Illuminate\Contracts\Http\Kernel as HttpKernelContract;
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

$kernel = $app->make(HttpKernelContract::class);
$getReq = Request::create('/login', 'GET');
$res = $kernel->handle($getReq);

$session = $getReq->session();
echo "started=".var_export($session->isStarted(), true)."\n";
echo "id=".$session->getId()."\n";
echo "driver=".config('session.driver')."\n";
echo "path=".config('session.files')."\n";
echo "handler=".get_class($session->getHandler())."\n";
$path = config('session.files').'/'.$session->getId();
echo "before_save exists=". (file_exists($path) ? 'yes' : 'no')."\n";

try {
    $session->save();
    echo "manual_save=ok\n";
} catch (Throwable $e) {
    echo "manual_save_err=".$e->getMessage()."\n".$e->getTraceAsString()."\n";
}
echo "after_manual exists=". (file_exists($path) ? 'yes' : 'no')." len=".(file_exists($path)?filesize($path):0)."\n";

$rawPath = storage_path('framework/sessions/_probe_put_test');
$n = file_put_contents($rawPath, 'hello', LOCK_EX);
echo "file_put_contents=$n exists=". (file_exists($rawPath) ? 'yes' : 'no')."\n";
@unlink($rawPath);

$files = app('files');
try {
    $files->put(storage_path('framework/sessions/_probe_fs_put'), 'world', true);
    echo "fs_put=ok exists=". (file_exists(storage_path('framework/sessions/_probe_fs_put')) ? 'yes' : 'no')."\n";
    @unlink(storage_path('framework/sessions/_probe_fs_put'));
} catch (Throwable $e) {
    echo "fs_put_err=".$e->getMessage()."\n";
}

$kernel->terminate($getReq, $res);
echo "after_terminate exists=". (file_exists($path) ? 'yes' : 'no')."\n";
