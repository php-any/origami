<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Pagination\LengthAwarePaginator;

$items = ['a', 'b', 'c'];
$paginator = new LengthAwarePaginator($items, 10, 3, 2, [
    'path' => '/posts',
    'pageName' => 'page',
]);

if ($paginator->total() !== 10) {
    echo "FAIL: total\n";
    exit(1);
}

if ($paginator->perPage() !== 3) {
    echo "FAIL: perPage\n";
    exit(1);
}

if ($paginator->currentPage() !== 2) {
    echo "FAIL: currentPage\n";
    exit(1);
}

if ($paginator->lastPage() !== 4) {
    echo "FAIL: lastPage expected 4\n";
    exit(1);
}

if ($paginator->count() !== 3) {
    echo "FAIL: count\n";
    exit(1);
}

if (!$paginator->hasMorePages()) {
    echo "FAIL: hasMorePages\n";
    exit(1);
}

$url = $paginator->url(3);
if (!is_string($url) || $url === '') {
    echo "FAIL: url\n";
    exit(1);
}

echo "PASS\n";
