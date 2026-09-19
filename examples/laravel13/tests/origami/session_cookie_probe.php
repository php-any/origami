<?php
$_SERVER['HTTP_HOST'] = '127.0.0.1';
$_SERVER['SERVER_NAME'] = '127.0.0.1';
$_SERVER['REQUEST_URI'] = '/';
$_SERVER['REQUEST_METHOD'] = 'GET';
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$cookie = null;
foreach (file(__DIR__.'/../../storage/origami-debug/e2e-cookies.txt') as $line) {
    $line = trim($line);
    if ($line === '' || (str_starts_with($line, '#') && !str_starts_with($line, '#HttpOnly_'))) {
        continue;
    }
    if (!str_contains($line, 'laravel-session')) {
        continue;
    }
    $parts = preg_split('/\s+/', $line);
    $cookie = $parts[count($parts) - 1] ?? null;
}
echo 'cookie_len='.strlen((string) $cookie)."\n";
try {
    $encrypter = $app->make('encrypter');
    $plain = $encrypter->decrypt($cookie, false);
    echo 'plain='.$plain."\n";
    $keys = $encrypter->getAllKeys();
    echo 'keys_type='.gettype($keys).' keys_count='.(is_countable($keys) ? count($keys) : -1)."\n";
    $id = Illuminate\Cookie\CookieValuePrefix::validate('laravel-session', $plain, $keys);
    echo 'validated_id='.var_export($id, true)."\n";
} catch (Throwable $e) {
    echo 'decrypt_ex='.$e->getMessage()."\n";
    exit;
}
$row = Illuminate\Support\Facades\DB::table('sessions')->where('id', $id)->first();
if (!$row) {
    echo "session_row=MISSING\n";
    exit;
}
$raw = base64_decode((string) $row->payload);
echo 'has_login='.(str_contains($raw, 'login_admin') ? 'yes' : 'no')."\n";
echo 'payload='.$raw."\n";

$session = $app->make('session.store');
$session->setId($id);
$session->start();
echo 'store_login='.var_export($session->get('login_admin_59ba36addc2b2f9401580f014c7f58ea4e30989d'), true)."\n";
echo 'guard_id='.var_export(Illuminate\Support\Facades\Auth::guard('admin')->id(), true)."\n";
echo 'guard_check='.(Illuminate\Support\Facades\Auth::guard('admin')->check() ? 'yes' : 'no')."\n";
$u = Illuminate\Support\Facades\Auth::guard('admin')->user();
echo 'guard_user='.($u ? $u->email : 'null')."\n";
