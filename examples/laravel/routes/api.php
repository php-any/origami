<?php

use App\Http\Controllers\ApiController;
use App\Http\Controllers\AuthController;
use App\Http\Controllers\PostApiController;
use App\Http\Controllers\UserController;
use App\Http\Middleware\Authenticate;
use App\Http\Middleware\LogRequest;
use Bootstrap\Routing\Route;

/*
|--------------------------------------------------------------------------
| API Routes（Route Facade / Net\Http\Router）
|--------------------------------------------------------------------------
|
| API 控制器不使用 #[GetMapping] 等注解，全部由本文件声明注册，
| 与 Web 侧的控制器注解路由共存。
|
*/

Route::group([
    'prefix' => 'api',
    'middleware' => [LogRequest::class],
], function () {
    Route::get('/health', [ApiController::class, 'health']);
    Route::post('/login', [AuthController::class, 'login']);

    Route::group(['middleware' => [Authenticate::class]], function () {
        Route::get('/me', [AuthController::class, 'me']);
        Route::get('/users', [UserController::class, 'index']);
        Route::get('/users/{id}', [UserController::class, 'show']);
        Route::get('/posts', [PostApiController::class, 'index']);
        Route::post('/posts', [PostApiController::class, 'store']);
        Route::get('/posts/{id}', [PostApiController::class, 'show']);
    });
});
