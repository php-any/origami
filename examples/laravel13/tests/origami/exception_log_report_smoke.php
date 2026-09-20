<?php

namespace tests\origami;

/**
 * Handler::report 必须把应上报的异常写入 storage/logs/laravel.log。
 * 依赖 with() 返回 shouldntReport 回调的 false。
 */

use Illuminate\Contracts\Debug\ExceptionHandler;
use RuntimeException;
use Symfony\Component\HttpKernel\Exception\NotFoundHttpException;

require __DIR__.'/../../vendor/autoload.php';

$app = require __DIR__.'/../../bootstrap/app.php';
$app->bootstrapWith([
    \Illuminate\Foundation\Bootstrap\LoadEnvironmentVariables::class,
    \Illuminate\Foundation\Bootstrap\LoadConfiguration::class,
    \Illuminate\Foundation\Bootstrap\HandleExceptions::class,
    \Illuminate\Foundation\Bootstrap\RegisterFacades::class,
    \Illuminate\Foundation\Bootstrap\SetRequestForConsole::class,
    \Illuminate\Foundation\Bootstrap\RegisterProviders::class,
    \Illuminate\Foundation\Bootstrap\BootProviders::class,
]);

$handler = $app->make(ExceptionHandler::class);
$runtime = new RuntimeException('origami-exception-should-report-to-laravel-log');
if (!$handler->shouldReport($runtime)) {
    \Log::fatal('RuntimeException 的 shouldReport 应为 true（with() 必须返回回调 false）');
}

$notFound = new NotFoundHttpException('missing');
if ($handler->shouldReport($notFound)) {
    \Log::fatal('NotFoundHttpException 属于 HttpException，shouldReport 应为 false');
}

$logPath = storage_path('logs/laravel.log');
$handler->report($runtime);
if (!is_file($logPath)) {
    \Log::fatal('report() 后应创建 '.$logPath);
}
$contents = file_get_contents($logPath);
if (!str_contains($contents, 'origami-exception-should-report-to-laravel-log')) {
    \Log::fatal('laravel.log 未写入 RuntimeException 消息: '.$contents);
}

\Log::info('exception log report smoke 通过');
