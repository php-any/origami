<?php

namespace App\Http\Middleware;

use Illuminate\Log\Events\MessageLogged;
use Net\Http\Request;
use Net\Http\Response;

/**
 * 请求日志中间件：控制台输出 + MessageLogged（供 Telescope LogWatcher 采集）。
 */
class LogRequest
{
    public function handle(Request $request, Response $response, callable $next): void
    {
        $start = microtime(true);
        $path = (string) $request->path();
        $method = (string) $request->method();

        $this->dispatchLog('info', '[HTTP] ' . $method . ' ' . $path);

        $next($request, $response);

        $ms = round((microtime(true) - $start) * 1000, 2);
        $this->dispatchLog('info', '[HTTP] 完成 ' . $path . ' (' . $ms . 'ms)');
    }

    private function dispatchLog(string $level, string $message, array $context = []): void
    {
        \Log::info($message);

        try {
            illuminate_container()['events']->dispatch(new MessageLogged($level, $message, $context));
        } catch (\Throwable $e) {
            // events 未就绪时仍保留控制台日志
        }
    }
}
