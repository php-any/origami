<?php

use Illuminate\Contracts\Http\Kernel as HttpKernelContract;
use Illuminate\Http\Request;
use Illuminate\Cookie\CookieValuePrefix;

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
$html = (string) $res->getContent();
preg_match('/data-csrf="([^"]*)"/', $html, $c);
$pageCsrf = $c[1] ?? '';
echo "during_get session_id=".session()->getId()." token=".session()->token()." page=$pageCsrf\n";
$kernel->terminate($getReq, $res);

$cookies = [];
foreach ($res->headers->getCookies() as $cookie) {
    $cookies[$cookie->getName()] = $cookie->getValue();
}
$raw = $cookies[config('session.cookie')];
$sid = CookieValuePrefix::remove(app('encrypter')->decrypt($raw, false));
$path = storage_path('framework/sessions/'.$sid);
echo "session_file_exists=". (file_exists($path) ? 'yes' : 'no')." path=$path\n";
if (file_exists($path)) {
    echo "session_file_len=".filesize($path)."\n";
    echo "session_file_has_token=". (str_contains(file_get_contents($path), '_token') ? 'yes' : 'no')."\n";
}

$req = Request::create('/x', 'GET', [], $cookies);
$app->make(\Illuminate\Cookie\Middleware\EncryptCookies::class)->handle($req, function ($request) use ($pageCsrf, $sid) {
    return app(\Illuminate\Session\Middleware\StartSession::class)->handle($request, function ($request) use ($pageCsrf, $sid) {
        echo "loaded_id=".$request->session()->getId()." expect=$sid same=".($request->session()->getId()===$sid?'yes':'no')."\n";
        echo "loaded_token=".$request->session()->token()." page=$pageCsrf match=".var_export(hash_equals($request->session()->token(), $pageCsrf), true)."\n";
        return response('ok');
    });
});
