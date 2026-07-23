<?php

/**
 * Eloquent + #[Table] Entity 双轨冒烟（需已 migrate 的 storage/laravel.db）。
 */

require dirname(__DIR__) . '/bootstrap/app.php';
bootstrap_app();
bootstrap_cli_container();

use App\Models\Entity\User as EntityUser;
use App\Models\Post;
use App\Models\User;
use App\Services\DatabaseManager;

app_make(DatabaseManager::class)->connect();

if (!class_exists(EntityUser::class)) {
    echo "FAIL: Entity User missing\n";
    exit(1);
}

$alice = User::query()->where('email', 'alice@example.com')->first();
if ($alice === null || $alice->name !== 'Alice') {
    echo "FAIL: Eloquent User query\n";
    exit(1);
}

if (Post::query()->count() < 1) {
    echo "FAIL: no posts\n";
    exit(1);
}

echo "PASS\n";
