<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Cache\ArrayStore;
use Illuminate\Cache\Repository;

/**
 * 不使用带 TTL 的 put/remember（会走 Carbon::now()，依赖 parent:: 等上游能力）。
 */
$cache = new Repository(new ArrayStore());

$cache->put('name', 'origami');
if ($cache->get('name') !== 'origami') {
    echo "FAIL: get after put\n";
    exit(1);
}

if (!$cache->has('name')) {
    echo "FAIL: has\n";
    exit(1);
}

$cache->forever('token', 'abc');
if ($cache->get('token') !== 'abc') {
    echo "FAIL: forever\n";
    exit(1);
}

$cache->forget('name');
if ($cache->has('name')) {
    echo "FAIL: forget\n";
    exit(1);
}

echo "PASS\n";
