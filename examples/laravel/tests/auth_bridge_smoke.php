<?php

/**
 * 阶段 2：Illuminate Auth Guard（api_token）冒烟。
 */

require dirname(__DIR__) . '/bootstrap/app.php';
bootstrap_app();
bootstrap_cli_container();

use App\Models\User;
use App\Services\DatabaseManager;
use Bootstrap\Auth\ApiTokenGuard;
use Illuminate\Database\Capsule\Manager as Capsule;

app_make(DatabaseManager::class)->connect();

$guard = auth()->guard();
if (!($guard instanceof ApiTokenGuard)) {
    echo "FAIL: default guard not ApiTokenGuard\n";
    exit(1);
}

$user = User::query()->where('email', 'alice@example.com')->first();
if ($user === null) {
    echo "FAIL: seed user missing (run migrate + db:seed)\n";
    exit(1);
}

$token = 'smoke-auth-' . bin2hex(random_bytes(8));
Capsule::table('api_tokens')->insert([
    'user_id' => $user->id,
    'token' => $token,
    'created_at' => '2026-07-22 00:00:00',
]);

auth_set_token($token);

if (!auth()->check()) {
    echo "FAIL: auth()->check()\n";
    exit(1);
}

$authed = auth()->user();
if ($authed === null) {
    echo "FAIL: auth()->user() null\n";
    exit(1);
}

if ($authed->id != $user->id) {
    echo "FAIL: auth()->user() id mismatch\n";
    echo "user=";
    var_export($user->id);
    echo " authed=";
    var_export($authed->id);
    echo "\n";
    exit(1);
}

auth_set_token('invalid-token');
if (!auth()->guest()) {
    echo "FAIL: guest after bad token\n";
    exit(1);
}

echo "PASS\n";
