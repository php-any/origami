<?php

// step11: symfony/http-kernel (完整请求生命周期)

require dirname(__DIR__) . '/bootstrap.php';

$kernel = symfony_create_kernel();

$response = symfony_handle($kernel, 'GET', '/');
step_check('HttpKernel GET / status', $response->getStatusCode() === 200, 'status=' . $response->getStatusCode());
$content = $response->getContent();
step_check('HttpKernel GET / content', str_contains($content ?? '', 'Origami Symfony'), substr($content ?? '', 0, 100));

$health = symfony_handle($kernel, 'GET', '/health');
step_check('HttpKernel GET /health', $health->getStatusCode() === 200);
step_check('HttpKernel health json', str_contains($health->getContent() ?? '', '"status"'));

step_check('step11_http_kernel', true, 'symfony/http-kernel');
