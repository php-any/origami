<?php

/**
 * Illuminate Auth 桥接（阶段 2：api Token Guard）
 */

use Bootstrap\Auth\ApiTokenGuard;
use Illuminate\Auth\AuthManager;
use Illuminate\Hashing\HashManager;

/** @var string|null 当前请求 Authorization token */
$GLOBALS['__origami_auth_token'] = null;

function auth_set_token(?string $token): void
{
    $GLOBALS['__origami_auth_token'] = $token ?? '';
    // 换请求后清掉已解析的 Guard 用户缓存
    if (function_exists('illuminate_auth')) {
        illuminate_auth()->forgetGuards();
    }
}

function auth_set_request($request): void
{
    $token = '';
    if (is_object($request) && method_exists($request, 'header')) {
        $token = (string) $request->header('Authorization', '');
    }
    auth_set_token($token);
}

function auth_current_token(): string
{
    return (string) ($GLOBALS['__origami_auth_token'] ?? '');
}

function illuminate_auth(): AuthManager
{
    return illuminate_container()->make('auth');
}

/**
 * @param  string|null  $guard
 * @return \Illuminate\Contracts\Auth\Guard|\Illuminate\Auth\AuthManager
 */
function auth($guard = null)
{
    $auth = illuminate_auth();
    if ($guard !== null) {
        return $auth->guard($guard);
    }

    return $auth;
}

function bootstrap_auth(): void
{
    $app = illuminate_container();

    if ($app->bound('auth')) {
        return;
    }

    if (!$app->bound('hash')) {
        $app->singleton('hash', function ($app) {
            return new HashManager($app);
        });
    }

    $app->singleton('auth', function ($app) {
        $auth = new AuthManager($app);
        $auth->extend('api_token', function () {
            return new ApiTokenGuard(static fn () => auth_current_token());
        });

        return $auth;
    });
}
