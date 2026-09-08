<?php

// 定位登录卡死：Auth::attempt / 邮箱校验是否会无限涨内存。
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

$log = __DIR__.'/_auth_attempt_progress.txt';
@unlink($log);
function ap($m) {
    global $log;
    file_put_contents($log, $m."\n", FILE_APPEND);
}

ap('boot_ok');
$v = validator(['email' => 'admin@example.com', 'password' => 'password'], [
    'email' => 'required|email',
    'password' => 'required',
]);
ap('validator_created');
if ($v->fails()) {
    ap('validate_fail '.json_encode($v->errors()->toArray()));
} else {
    ap('validate_ok');
}
ap('attempt_start');
$ok = auth('admin')->attempt(['email' => 'admin@example.com', 'password' => 'password']);
ap('attempt='.var_export($ok, true));
if ($ok) {
    $u = auth('admin')->user();
    ap('user='.($u ? $u->email : 'null').' active='.var_export($u?->is_active, true));
    ap('redirect_start');
    $url = redirect()->route('admin.dashboard')->getTargetUrl();
    ap('redirect='.$url);
}
ap('done');
echo file_get_contents($log);
