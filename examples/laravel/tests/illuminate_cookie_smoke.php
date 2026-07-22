<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Cookie\CookieJar;
use Symfony\Component\HttpFoundation\Cookie;

/**
 * minutes=0 为会话 Cookie，避免 Carbon::now()（parent:: 递归问题）。
 */
$jar = new CookieJar();
$jar->queue($jar->make('token', 'abc123', 0));

$queued = $jar->getQueuedCookies();
if (count($queued) !== 1) {
    echo "FAIL: queue count\n";
    exit(1);
}

/** @var Cookie $cookie */
$cookie = $queued[0];
if ($cookie->getName() !== 'token' || $cookie->getValue() !== 'abc123') {
    echo "FAIL: cookie name/value\n";
    var_export($cookie);
    echo "\n";
    exit(1);
}

$jar->unqueue('token');
if (count($jar->getQueuedCookies()) !== 0) {
    echo "FAIL: unqueue\n";
    exit(1);
}

echo "PASS\n";
