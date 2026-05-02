<?php

namespace Gateway;

use Net\Http\Server;
use Net\Http\app;

/**
 * Gateway Example
 *
 * This example demonstrates a powerful API Gateway that routes requests
 * to different applications based on the Host header.
 *
 * It uses the `app()` function to dynamically load and execute
 * different applications in isolated environments.
 */

// Define the routing table
// Host -> { path, func }
$routes = [
    "demo.local" => [
        "path" => "apps/demo/src/main.php",
        "func" => "DemoApp\\main"
    ],
    "spring.local" => [
        "path" => "../spring/src/main.php",
        "func" => "App\\main"
    ],
    "default" => [
        "path" => "apps/demo/src/main.php",
        "func" => "DemoApp\\main"
    ]
];

$port = 8080;
$server = new Server("0.0.0.0", port: $port);

// Collect route keys for logging
$routeKeys = [];
foreach ($routes as $k => $v) {
    $routeKeys[] = $k;
}
Log::info("Starting Gateway on port {$port}...");
Log::info("Routes configured: " . implode(", ", $routeKeys));

// Middleware for logging
$server->middleware(function ($req, $res, $next) {
    Log::info("[Gateway] Request: " . $req->method() . " " . $req->fullUrl());
    $next($req, $res);
});

// Main handler
$server->any(function ($req, $res) {
    // Get Host header
    $host = $req->header("Host");
    if ($host == null) {
        $host = "default";
    }

    // Remove port from host if present (e.g. localhost:8080 -> localhost)
    if (strpos($host, ":") !== false) {
        $parts = explode(":", $host);
        $host = $parts[0];
    }

    Log::debug("Routing for host: " . $host);

    // Lookup route
    $route = $routes[$host] ?? null;
    if ($route == null) {
        // Fallback to default if configured
        if (isset($routes["default"])) {
            $route = $routes["default"];
        } else {
            $res->writeHeader(404);
            $res->json([
                "error" => "Service not found",
                "host" => $host
            ]);
            return;
        }
    }

    // Dispatch to the application
    try {
        // app() function loads the file in a new/temp VM and calls the function
        // It provides isolation between requests/apps
        app($req, $res, $route["path"], $route["func"]);
    } catch (\Exception $e) {
        Log::error("Gateway Dispatch Error: " . $e->getMessage());
        $res->writeHeader(500);
        $res->json([
            "error" => "Gateway Internal Error",
            "message" => $e->getMessage()
        ]);
    }
});

$server->run();
