<?php


return [
    'default' => env('DB_CONNECTION', 'sqlite'),
    'connections' => [
        'sqlite' => [
            'driver' => 'sqlite',
            'database' => realpath(env('DB_DATABASE', __DIR__ . '/../storage/laravel.db')),
        ],
    ],
];
