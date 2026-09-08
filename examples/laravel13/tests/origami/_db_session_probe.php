<?php

use Illuminate\Contracts\Http\Kernel as HttpKernelContract;
use Illuminate\Http\Request;
use Illuminate\Cookie\CookieValuePrefix;
use Illuminate\Support\Facades\DB;

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
$sid = session()->getId();
$token = session()->token();
echo "during_get id=$sid token=$token page=$pageCsrf\n";

$row = DB::table(config('session.table', 'sessions'))->where('id', $sid)->first();
echo "db_row_after_handle=". ($row ? 'yes' : 'no')."\n";
if ($row) {
    echo "db_payload_len=".strlen($row->payload)."\n";
    echo "db_has_token=". (str_contains(base64_decode($row->payload), '_token') || str_contains($row->payload, '_token') ? 'yes' : 'no')."\n";
}

$kernel->terminate($getReq, $res);

$row2 = DB::table(config('session.table', 'sessions'))->where('id', $sid)->first();
echo "db_row_after_terminate=". ($row2 ? 'yes' : 'no')."\n";

$cookies = [];
foreach ($res->headers->getCookies() as $cookie) {
    $cookies[$cookie->getName()] = $cookie->getValue();
}
$raw = $cookies[config('session.cookie')] ?? '';
echo "cookie_len=".strlen($raw)."\n";
try {
    $dec = CookieValuePrefix::remove(app('encrypter')->decrypt($raw, false));
    echo "decrypt_sid=$dec match=".($dec===$sid?'yes':'no')."\n";
} catch (Throwable $e) {
    echo "decrypt_err=".$e->getMessage()."\n";
}

$req = Request::create('/x', 'GET', [], $cookies);
$app->make(\Illuminate\Cookie\Middleware\EncryptCookies::class)->handle($req, function ($request) use ($pageCsrf, $sid) {
    return app(\Illuminate\Session\Middleware\StartSession::class)->handle($request, function ($request) use ($pageCsrf, $sid) {
        echo "loaded_id=".$request->session()->getId()." expect=$sid same=".($request->session()->getId()===$sid?'yes':'no')."\n";
        echo "loaded_token=".$request->session()->token()." page=$pageCsrf match=".var_export(hash_equals($request->session()->token(), $pageCsrf), true)."\n";
        $row = DB::table(config('session.table', 'sessions'))->where('id', $sid)->first();
        echo "db_at_load=". ($row ? 'yes payload_len='.strlen($row->payload) : 'no')."\n";
        return response('ok');
    });
});
