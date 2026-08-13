<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Session\ArraySessionHandler;
use Illuminate\Session\Store;

/**
 * 不使用 flash（依赖 Arr::set 点号键，Origami 尚未完整支持）。
 */
$handler = new ArraySessionHandler();
$session = new Store('origami_smoke', $handler);
$session->start();

$session->put('user', 'alice');
if ($session->get('user') !== 'alice') {
    echo "FAIL: put/get\n";
    exit(1);
}

if (!$session->has('user')) {
    echo "FAIL: has\n";
    exit(1);
}

$session->forget('user');
if ($session->has('user')) {
    echo "FAIL: forget\n";
    exit(1);
}

$id = $session->getId();
if (!is_string($id) || $id === '') {
    echo "FAIL: session id\n";
    exit(1);
}

echo "PASS\n";
