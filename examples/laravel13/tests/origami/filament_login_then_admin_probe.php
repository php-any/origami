<?php
/**
 * Livewire authenticate 后检查 session 是否已登录，并跟进 GET /admin。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$kernel = $app->make(Illuminate\Contracts\Http\Kernel::class);

$get = $kernel->handle(Illuminate\Http\Request::create('http://127.0.0.1:8000/admin/login', 'GET'));
$html = (string) $get->getContent();
preg_match('/wire:snapshot="([^"]+)"/', $html, $s);
preg_match('/data-csrf="([^"]*)"/', $html, $c);
preg_match('/data-update-uri="([^"]+)"/', $html, $u);
$snapshot = html_entity_decode($s[1] ?? '', ENT_QUOTES);
$csrf = $c[1] ?? '';
$updateUri = html_entity_decode($u[1] ?? '', ENT_QUOTES);
$cookies = [];
foreach ($get->headers->getCookies() as $cookie) {
    $cookies[$cookie->getName()] = $cookie->getValue();
}

$payload = json_encode([
    'components' => [[
        'snapshot' => $snapshot,
        'updates' => [
            'data.email' => 'admin@example.com',
            'data.password' => 'password',
        ],
        'calls' => [
            ['method' => 'authenticate', 'params' => []],
        ],
    ]],
]);
$server = [
    'CONTENT_TYPE' => 'application/json',
    'HTTP_ACCEPT' => 'application/json',
    'HTTP_X_LIVEWIRE' => 'true',
    'HTTP_X_CSRF_TOKEN' => $csrf,
    'CONTENT_LENGTH' => (string) strlen($payload),
];
$post = $kernel->handle(Illuminate\Http\Request::create($updateUri, 'POST', [], $cookies, [], $server, $payload));
$body = (string) $post->getContent();
$json = json_decode($body, true);
echo "post_status=".$post->getStatusCode()."\n";
echo "redirect=".var_export($json['components'][0]['effects']['redirect'] ?? null, true)."\n";

foreach ($post->headers->getCookies() as $cookie) {
    $cookies[$cookie->getName()] = $cookie->getValue();
}

// 同进程查 guard（Livewire 请求后）
echo "guard_after_post=".(Filament\Facades\Filament::auth()->check() ? 'yes' : 'no')."\n";
echo "session_has=".(session()->has('login_admin_'.sha1('Illuminate\\Auth\\SessionGuard')) || session()->all() ? 'dump' : 'no')."\n";
echo "session_keys=".implode(',', array_keys(session()->all()))."\n";

// GET /admin with cookies
$adminReq = Illuminate\Http\Request::create('http://127.0.0.1:8000/admin', 'GET', [], $cookies);
try {
    $admin = $kernel->handle($adminReq);
    echo "admin_status=".$admin->getStatusCode()."\n";
    echo "admin_loc=".(string)$admin->headers->get('Location')."\n";
    $ac = (string)$admin->getContent();
    echo "admin_len=".strlen($ac)."\n";
    echo "admin_markers=".(str_contains($ac, 'Origami Admin') || str_contains($ac, 'Dashboard') || str_contains($ac, '仪表盘') ? 'yes' : 'no')."\n";
    if ($admin->getStatusCode() >= 400) {
        echo "admin_head=".substr(preg_replace('/\s+/', ' ', strip_tags($ac)), 0, 300)."\n";
    }
} catch (Throwable $e) {
    echo "admin EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
}
echo "DONE\n";
