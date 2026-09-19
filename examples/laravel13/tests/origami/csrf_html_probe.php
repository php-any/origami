<?php
/**
 * 在 Laravel HTTP Kernel 里打印 csrf_token 与登录页 HTML。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$kernel = $app->make(Illuminate\Contracts\Http\Kernel::class);

$req = Illuminate\Http\Request::create('http://127.0.0.1:18086/admin/login', 'GET');
$get = $kernel->handle($req);
$html = (string) $get->getContent();
echo "status=".$get->getStatusCode()."\n";
echo "has_session_store=".(app()->has('session.store') ? 'yes' : 'no')."\n";
echo "session_class=".get_class(app('session'))."\n";
try {
    $tok = csrf_token();
    echo "csrf_token=".var_export($tok, true)." len=".strlen((string) $tok)."\n";
} catch (Throwable $e) {
    echo "csrf_token_ex=".$e->getMessage()."\n";
}
try {
    echo "session_token=".var_export(app('session')->token(), true)."\n";
} catch (Throwable $e) {
    echo "session_token_ex=".$e->getMessage()."\n";
}
try {
    echo "request_session_token=".var_export($req->session()->token(), true)."\n";
} catch (Throwable $e) {
    echo "request_session_ex=".$e->getMessage()."\n";
}
preg_match('/name="csrf-token" content="([^"]*)"/', $html, $m);
echo "html_csrf=".var_export($m[1] ?? 'MISSING', true)."\n";
preg_match('/data-csrf="([^"]*)"/', $html, $d);
echo "data_csrf=".var_export($d[1] ?? 'MISSING', true)."\n";
