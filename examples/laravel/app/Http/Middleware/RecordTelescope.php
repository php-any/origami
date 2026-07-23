<?php

namespace App\Http\Middleware;

use Laravel\Telescope\IncomingEntry;
use Laravel\Telescope\Telescope;
use Net\Http\Request;
use Net\Http\Response;

/**
 * 请求结束时写入 Telescope request/log 并 store（替代 Http Kernel terminating）。
 */
class RecordTelescope
{
    public function handle(Request $request, Response $response, callable $next): void
    {
        bootstrap_telescope_http();

        $path = ltrim((string) $request->path(), '/');
        $skip = $path === 'telescope'
            || str_starts_with($path, 'telescope/')
            || str_starts_with($path, 'vendor/telescope');

        $paused = false;
        try {
            $paused = (bool) cache('telescope:pause-recording');
        } catch (\Throwable $e) {
        }

        $record = !$skip && config('telescope.enabled', true) && !$paused;
        $start = microtime(true);

        if ($record) {
            Telescope::startRecording(false);
        }

        try {
            $next($request, $response);
        } finally {
            if (!$record || !Telescope::isRecording()) {
                return;
            }

            $ms = (int) floor((microtime(true) - $start) * 1000);
            $uri = '/' . $path;
            if ($uri === '/') {
                $uri = '/';
            }

            try {
                Telescope::recordRequest(IncomingEntry::make([
                    'ip_address' => method_exists($request, 'ip') ? (string) $request->ip() : '',
                    'uri' => $uri,
                    'method' => (string) $request->method(),
                    'controller_action' => '',
                    'middleware' => [],
                    'headers' => [],
                    'payload' => [],
                    'session' => [],
                    'response_status' => 200,
                    'response' => [],
                    'duration' => $ms,
                    'memory' => 0,
                ]));

                Telescope::recordLog(IncomingEntry::make([
                    'level' => 'info',
                    'message' => '[HTTP] ' . $request->method() . ' ' . $uri . ' (' . $ms . 'ms)',
                    'context' => [],
                ]));

                Telescope::store(telescope_entries_repository());
            } catch (\Throwable $e) {
                \Log::error('Telescope store failed: ' . $e->getMessage());
            }
        }
    }
}
