<?php

use Illuminate\Pagination\LengthAwarePaginator;
use Illuminate\Pagination\Paginator;

if (!class_exists(LengthAwarePaginator::class, false)) {
    echo "skip: LengthAwarePaginator 未注册\n";
    return;
}

$la = new LengthAwarePaginator([1, 2, 3], 10, 3, 2, ['path' => '/items']);
if ($la->total() !== 10 || $la->perPage() !== 3 || $la->currentPage() !== 2) {
    echo "FAIL LengthAwarePaginator basics\n";
    exit(1);
}
if ($la->lastPage() !== 4 || !$la->hasMorePages()) {
    echo "FAIL LengthAwarePaginator lastPage/hasMorePages\n";
    exit(1);
}
if ($la->links() !== '') {
    echo "FAIL links stub\n";
    exit(1);
}
$arr = $la->toArray();
if (($arr['total'] ?? null) !== 10 || ($arr['current_page'] ?? null) !== 2) {
    echo "FAIL toArray " . json_encode($arr) . "\n";
    exit(1);
}

if (class_exists(Paginator::class, false)) {
    $p = new Paginator([1, 2, 3, 4, 5], 2, 1);
    if ($p->count() !== 2 || !$p->hasMorePages()) {
        echo "FAIL Paginator slice/hasMore\n";
        exit(1);
    }
}

echo "OK illuminate_pagination_test\n";
