<?php

use App\Http\Controllers\Telescope\ApiController as TelescopeApiController;
use App\Http\Controllers\Telescope\DashboardController;
use App\Http\Controllers\HomeController;
use App\Http\Controllers\PostController;
use App\Http\Middleware\LogRequest;
use App\Http\Middleware\RecordTelescope;
use Bootstrap\Routing\Route;

/*
|--------------------------------------------------------------------------
| Web Routes
|--------------------------------------------------------------------------
*/

Route::group(['middleware' => [LogRequest::class, RecordTelescope::class]], function () {
    Route::get('/', [HomeController::class, 'index']);
    Route::get('/posts', [PostController::class, 'index']);
    Route::get('/posts/{id}', [PostController::class, 'show']);
});

/*
|--------------------------------------------------------------------------
| Telescope Dashboard + API（Origami 分发，非官方 Illuminate Route）
|--------------------------------------------------------------------------
*/
Route::group(['prefix' => 'telescope'], function () {
    // API（须先于 SPA catch-all）
    Route::post('/telescope-api/requests', [TelescopeApiController::class, 'requestsIndex']);
    Route::get('/telescope-api/requests/{id}', [TelescopeApiController::class, 'requestsShow']);
    Route::post('/telescope-api/logs', [TelescopeApiController::class, 'logsIndex']);
    Route::get('/telescope-api/logs/{id}', [TelescopeApiController::class, 'logsShow']);
    Route::post('/telescope-api/queries', [TelescopeApiController::class, 'queriesIndex']);
    Route::get('/telescope-api/queries/{id}', [TelescopeApiController::class, 'queriesShow']);
    Route::post('/telescope-api/exceptions', [TelescopeApiController::class, 'exceptionsIndex']);
    Route::get('/telescope-api/exceptions/{id}', [TelescopeApiController::class, 'exceptionsShow']);

    foreach ([
        'mail', 'dumps', 'notifications', 'jobs', 'batches', 'events', 'gates',
        'cache', 'models', 'views', 'commands', 'schedule', 'redis', 'client-requests',
    ] as $type) {
        Route::post('/telescope-api/' . $type, [TelescopeApiController::class, 'emptyIndex']);
        Route::get('/telescope-api/' . $type . '/{id}', [TelescopeApiController::class, 'emptyShow']);
    }

    Route::get('/telescope-api/monitored-tags', [TelescopeApiController::class, 'monitoredTags']);
    Route::post('/telescope-api/monitored-tags/', [TelescopeApiController::class, 'monitoredTagsStore']);
    Route::post('/telescope-api/monitored-tags/delete', [TelescopeApiController::class, 'monitoredTagsDelete']);
    Route::post('/telescope-api/toggle-recording', [TelescopeApiController::class, 'toggleRecording']);
    Route::delete('/telescope-api/entries', [TelescopeApiController::class, 'clearEntries']);

    Route::get('/', [DashboardController::class, 'index']);
    Route::get('/{view...}', [DashboardController::class, 'index']);
});
