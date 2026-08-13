<?php

require dirname(__DIR__) . '/bootstrap/app.php';
bootstrap_app();

$a = [
    'controller_action' => optional(null)->getActionName(),
    'middleware' => array_values(optional(null)->gatherMiddleware() ?? []),
];

if (count($a) !== 2) {
    Log::fatal('expected 2 keys, got ' . count($a));
}
if ($a['controller_action'] !== null) {
    Log::fatal('expected null controller_action');
}
if ($a['middleware'] !== []) {
    Log::fatal('expected empty middleware');
}

Log::info('laravel optional null array ok');
