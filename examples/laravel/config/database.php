<?php


return [
    'default' => env('DB_CONNECTION', 'sqlite'),
    'connections' => [
        'sqlite' => [
            'driver' => 'sqlite',
            'database' => env('DB_DATABASE', base_path('storage/laravel.db')),
            'prefix' => '',
        ],
    ],
];
