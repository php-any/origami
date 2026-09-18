<?php
/**
 * Livewire update 全程：监听 dehydrate 时的 redirect store / effects。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use Livewire\Component;

$log = [];
\Livewire\on('dehydrate', function ($component, $context) use (&$log) {
    $to = \Livewire\store($component)->get('redirect');
    $log[] = [
        'class' => get_class($component),
        'store_redirect' => $to,
        'effects_before' => array_keys($context->effects ?? []),
        'redirect_app' => get_class(app('redirect')),
        'is_lw_req' => app(\Livewire\Mechanisms\HandleRequests\HandleRequests::class)->isLivewireRequest(),
    ];
});

\Livewire\on('call', function ($component, $method, $params, $componentContext, $returnEarly, $metadata) use (&$log) {
    $log[] = [
        'event' => 'call_before',
        'method' => $method,
        'redirect_app' => get_class(app('redirect')),
    ];
    return function ($return) use (&$log, $component, $method) {
        $log[] = [
            'event' => 'call_after',
            'method' => $method,
            'return_type' => is_object($return) ? get_class($return) : gettype($return),
            'is_responsable' => is_object($return) && $return instanceof \Illuminate\Contracts\Support\Responsable,
            'store_redirect' => \Livewire\store($component)->get('redirect'),
            'redirect_app' => get_class(app('redirect')),
        ];
    };
});

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
], JSON_UNESCAPED_UNICODE);
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
echo "redirect_effect=".var_export($json['components'][0]['effects']['redirect'] ?? null, true)."\n";
echo "effect_keys=".implode(',', array_keys($json['components'][0]['effects'] ?? []))."\n";
echo "log=".json_encode($log, JSON_UNESCAPED_UNICODE|JSON_PRETTY_PRINT)."\n";
echo "DONE\n";
