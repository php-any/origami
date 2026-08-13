<?php

namespace App\Http\Middleware;

use Net\Http\Request;
use Net\Http\Response;

/**
 * 预留中间件：当前不注入任何 Telescope 专用录制逻辑。
 */
class RecordTelescope
{
    public function handle(Request $request, Response $response, callable $next): void
    {
        $next($request, $response);
    }
}
