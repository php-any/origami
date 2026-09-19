<?php
$_SERVER['HTTP_HOST'] = '127.0.0.1';
$_SERVER['SERVER_NAME'] = '127.0.0.1';
$_SERVER['REQUEST_URI'] = '/admin';
$_SERVER['REQUEST_METHOD'] = 'GET';
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$kernel = $app->make(Illuminate\Contracts\Http\Kernel::class);
if (method_exists($kernel, 'bootstrap')) {
    $kernel->bootstrap();
}

$cookie = null;
foreach (file(__DIR__.'/../../storage/origami-debug/e2e-cookies.txt') as $line) {
    $line = trim($line);
    if (!str_contains($line, 'laravel-session')) {
        continue;
    }
    $parts = preg_split('/\s+/', $line);
    $cookie = $parts[count($parts) - 1] ?? null;
}

$req = Illuminate\Http\Request::create('http://127.0.0.1:18086/admin', 'GET', [], [
    'laravel-session' => $cookie,
]);
echo 'cookie_all_keys='.implode(',', array_keys($req->cookies->all()))."\n";
$n = 0;
foreach ($req->cookies as $k => $v) {
    echo 'foreach k='.$k.' vlen='.strlen((string) $v)."\n";
    $n++;
}
echo "foreach_n=$n\n";
try {
    $mw = $app->make(Illuminate\Cookie\Middleware\EncryptCookies::class);
    $rp = new ReflectionProperty($mw, 'encrypter');
    $rp->setAccessible(true);
    $mwEnc = $rp->getValue($mw);
    echo 'mw_decrypt='.$mwEnc->decrypt($cookie, false)."\n";
$req2 = Illuminate\Http\Request::create('http://127.0.0.1:18086/admin', 'GET', [], [
    'laravel-session' => $cookie,
]);
$mwResp = $mw->handle($req2, function ($r) {
    echo 'after_encryptcookies='.var_export($r->cookies->get('laravel-session'), true)."\n";
    return response('ok');
});
} catch (Throwable $e) {
    echo 'mw_decrypt_ex='.$e->getMessage()."\n";
}
$res = $kernel->handle($req);
try {
    echo 'has_session='.($req->hasSession() ? 'yes' : 'no')."\n";
} catch (Throwable $e) {
    echo 'has_session_ex='.$e->getMessage()."\n";
}
echo 'facade_check='.(Illuminate\Support\Facades\Auth::guard('admin')->check() ? 'yes' : 'no')."\n";
try {
    echo 'facade_session_id='.Illuminate\Support\Facades\Session::getId()."\n";
    echo 'facade_login='.var_export(Illuminate\Support\Facades\Session::get('login_admin_59ba36addc2b2f9401580f014c7f58ea4e30989d'), true)."\n";
} catch (Throwable $e) {
    echo 'facade_sess_ex='.$e->getMessage()."\n";
}
echo 'status='.$res->getStatusCode()."\n";
echo 'loc='.(string) $res->headers->get('Location')."\n";
echo 'len='.strlen((string) $res->getContent())."\n";
$c = (string) $res->getContent();
echo 'markers='.(str_contains($c, '仪表盘') || str_contains($c, 'Origami Admin') || str_contains($c, 'Dashboard') ? 'yes' : 'no')."\n";
