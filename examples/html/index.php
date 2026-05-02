use Net\Http\Server;
use Net\Http\app;

$server = new Server("0.0.0.0", port: 8080);

// 简单的日志中间件
$server->middleware(function ($request, $response, $next) {
    $method = $request->method();
    $path = $request->path();
    Log::info("HTTP " . $method . " " . $path);
    $next($request, $response);
});

$server->any(function ($request, $response) {
    app($request, $response, "./main.php");
});

Log::info("HTML 服务启动在: http://127.0.0.1:8080");
$server->run();

