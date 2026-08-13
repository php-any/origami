<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Events\Dispatcher;

$events = new Dispatcher();
$box = new stdClass();
$box->seen = [];

// 使用 wildcard 签名 ($event, $payload)，避免依赖 argument unpacking (...$payload)
// 当前 Origami 对 $fn(...array_values($payload)) 会把整个数组当作第一个参数。
$events->listen('user.*', function ($event, $payload) use ($box) {
    $box->seen[] = [$event, $payload[0] ?? null];
});

$events->dispatch('user.created', ['alice']);

if (count($box->seen) !== 1 || $box->seen[0][0] !== 'user.created' || $box->seen[0][1] !== 'alice') {
    echo "FAIL: listener not invoked correctly\n";
    exit(1);
}

echo "PASS\n";
