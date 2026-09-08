<?php

use Illuminate\Contracts\Http\Kernel as HttpKernelContract;
use Illuminate\Http\Request;
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
$pageToken = session()->token();
$sid = session()->getId();
echo "page_token=$pageToken sid=$sid\n";
$store = session()->driver();
echo "serialization=". (function ($s) {
    $r = new ReflectionProperty($s, 'serialization');
    $r->setAccessible(true);
    return $r->getValue($s);
})($store)."\n";
$kernel->terminate($getReq, $res);

$row = DB::table('sessions')->where('id', $sid)->first();
$payload = base64_decode($row->payload);
echo "raw_payload=$payload\n";
$json = json_decode($payload, true);
echo "json_decode_ok=". (is_array($json) ? 'yes' : 'no')." json_token=".($json['_token'] ?? 'MISSING')." json_err=".json_last_error_msg()."\n";

$cookies = [];
foreach ($res->headers->getCookies() as $cookie) {
    $cookies[$cookie->getName()] = $cookie->getValue();
}

$req = Request::create('/x', 'GET', [], $cookies);
$app->make(\Illuminate\Cookie\Middleware\EncryptCookies::class)->handle($req, function ($request) use ($pageToken, $sid, $payload) {
    $handler = app('session')->getHandler();
    // Fresh store like StartSession would create
    $session = app('session')->driver();
    echo "fresh_driver_id_before_start=".$session->getId()."\n";

    return app(\Illuminate\Session\Middleware\StartSession::class)->handle($request, function ($request) use ($pageToken, $sid, $payload) {
        $session = $request->session();
        echo "started_id=".$session->getId()."\n";
        echo "started_token=".$session->token()."\n";
        echo "page_token=$pageToken match=".var_export(hash_equals($session->token(), $pageToken), true)."\n";
        echo "all=".json_encode($session->all())."\n";

        $handler = $session->getHandler();
        $read = $handler->read($sid);
        echo "direct_read_len=".strlen((string)$read)." eq_payload=".(($read === $payload)?'yes':'no')."\n";
        echo "direct_read=$read\n";
        $decoded = json_decode($read, true);
        echo "direct_json_token=".($decoded['_token'] ?? 'MISSING')."\n";

        // exists flag?
        $ref = new ReflectionProperty($handler, 'exists');
        $ref->setAccessible(true);
        echo "handler_exists=".var_export($ref->getValue($handler), true)."\n";

        return response('ok');
    });
});
