<?php

return [
    'defaults' => [
        'guard' => 'api',
        'provider' => 'users',
    ],

    'guards' => [
        'api' => [
            'driver' => 'api_token',
            'provider' => 'users',
        ],
    ],

    'providers' => [
        'users' => [
            'driver' => 'eloquent',
            'model' => App\Models\User::class,
        ],
    ],
];
