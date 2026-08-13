<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Auth\Access\Gate;
use Illuminate\Auth\GenericUser;
use Illuminate\Container\Container;

$user = new GenericUser([
    'id' => 42,
    'password' => 'secret',
    'remember_token' => 'token',
]);

if ($user->getAuthIdentifier() !== 42) {
    echo "FAIL: GenericUser id\n";
    exit(1);
}

if ($user->getAuthPassword() !== 'secret') {
    echo "FAIL: GenericUser password\n";
    exit(1);
}

$gate = new Gate(new Container(), fn () => $user);
$gate->define('smoke', fn ($authUser) => $authUser !== null && $authUser->getAuthIdentifier() === 42);

if (!$gate->allows('smoke')) {
    echo "FAIL: gate allows\n";
    exit(1);
}

if ($gate->denies('smoke')) {
    echo "FAIL: gate denies\n";
    exit(1);
}

echo "PASS\n";
