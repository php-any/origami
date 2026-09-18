<?php
/**
 * Filament /admin/login：GET 页面 + Livewire update 调用 authenticate。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$kernel = $app->make(Illuminate\Contracts\Http\Kernel::class);

echo "GET /admin/login...\n";
$get = $kernel->handle(Illuminate\Http\Request::create('http://127.0.0.1:8000/admin/login', 'GET'));
$html = (string) $get->getContent();
echo "status=".$get->getStatusCode()." len=".strlen($html)."\n";

if ($get->getStatusCode() !== 200) {
    echo "FAIL: login page status\n";
    exit(1);
}

// Livewire snapshot / csrf / update uri
preg_match('/wire:snapshot="([^"]+)"/', $html, $s);
preg_match('/data-csrf="([^"]*)"/', $html, $c);
preg_match('/data-update-uri="([^"]+)"/', $html, $u);
preg_match('/wire:id="([^"]+)"/', $html, $id);

$snapshotRaw = html_entity_decode($s[1] ?? '', ENT_QUOTES);
$csrf = $c[1] ?? '';
$updateUri = html_entity_decode($u[1] ?? '', ENT_QUOTES);
$wireId = $id[1] ?? '';

echo "has_snapshot=".($snapshotRaw !== '' ? 'yes' : 'no')."\n";
echo "has_csrf=".($csrf !== '' ? 'yes' : 'no')."\n";
echo "updateUri=".($updateUri !== '' ? $updateUri : 'MISSING')."\n";
echo "wireId=".($wireId !== '' ? $wireId : 'MISSING')."\n";

if ($snapshotRaw === '' || $updateUri === '') {
    // dump hints
    echo "has_livewire_js=".(str_contains($html, 'livewire') ? 'yes' : 'no')."\n";
    echo "has_authenticate=".(str_contains($html, 'authenticate') ? 'yes' : 'no')."\n";
    echo "has_wire_submit=".(str_contains($html, 'wire:submit') ? 'yes' : 'no')."\n";
    echo "FAIL: missing snapshot/updateUri\n";
    file_put_contents(__DIR__.'/_filament_login_page.html', $html);
    exit(1);
}

$cookies = [];
foreach ($get->headers->getCookies() as $cookie) {
    $cookies[$cookie->getName()] = $cookie->getValue();
}

// Livewire v3+ payload：snapshot 常为 JSON 字符串
$snapshot = $snapshotRaw;
$decoded = json_decode($snapshotRaw, true);
if (is_array($decoded)) {
    // 已是对象则保持字符串传给 components.snapshot（Livewire 期望字符串）
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
], JSON_UNESCAPED_UNICODE);

$server = [
    'CONTENT_TYPE' => 'application/json',
    'HTTP_ACCEPT' => 'application/json',
    'HTTP_X_LIVEWIRE' => 'true',
    'HTTP_X_CSRF_TOKEN' => $csrf,
    'HTTP_X_XSRF_TOKEN' => $cookies['XSRF-TOKEN'] ?? '',
    'CONTENT_LENGTH' => (string) strlen($payload),
];

echo "POST Livewire authenticate...\n";
try {
    $postReq = Illuminate\Http\Request::create($updateUri, 'POST', [], $cookies, [], $server, $payload);
    $post = $kernel->handle($postReq);
} catch (Throwable $e) {
    echo "EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
    exit(1);
}

$body = (string) $post->getContent();
echo "post_status=".$post->getStatusCode()." len=".strlen($body)."\n";
file_put_contents(__DIR__.'/_filament_login_post.json', $body);

if ($post->getStatusCode() >= 500) {
    echo "FAIL body: ".substr($body, 0, 800)."\n";
    exit(1);
}

$json = json_decode($body, true);
$redirect = $json['components'][0]['effects']['redirect'] ?? null;
$effects = $json['components'][0]['effects'] ?? null;
echo "redirect=".var_export($redirect, true)."\n";
if (isset($json['components'][0]['effects']['html'])) {
    echo "has_html_effect=yes\n";
}
if (isset($json['message'])) {
    echo "message=".$json['message']."\n";
}
if (isset($json['exception'])) {
    echo "exception=".$json['exception']."\n";
}

// 打印 effects 键便于诊断
if (is_array($effects)) {
    echo "effect_keys=".implode(',', array_keys($effects))."\n";
} else {
    echo "no_effects top_keys=".(is_array($json) ? implode(',', array_keys($json)) : gettype($json))."\n";
    echo "body_head=".substr($body, 0, 500)."\n";
}

echo "DONE\n";
