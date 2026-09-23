<?php

/**
 * JsonResponse / RedirectResponse 原生实现冒烟。未 AddClass 时跳过。
 */
if (!class_exists('Illuminate\\Http\\JsonResponse', false)) {
    echo "skip: JsonResponse 原生类未注册\n";
    return;
}

$json = new \Illuminate\Http\JsonResponse(['ok' => true, 'n' => 1], 201);
$content = $json->getContent();
$decoded = json_decode($content, true);
if (!is_array($decoded) || ($decoded['ok'] ?? null) !== true) {
    echo "FAIL JsonResponse content " . var_export($content, true) . "\n";
    exit(1);
}
if ($json->status() !== 201) {
    echo "FAIL JsonResponse status\n";
    exit(1);
}

$redirect = new \Illuminate\Http\RedirectResponse('http://example.com/login', 302);
if ($redirect->getTargetUrl() !== 'http://example.com/login') {
    echo "FAIL RedirectResponse target\n";
    exit(1);
}

echo "OK illuminate_http_response_test\n";
