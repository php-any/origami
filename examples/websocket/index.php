<?php

use Net\Http\Server;
use Net\Websocket\upgrade;

/**
 * WebSocket Echo 示例
 *
 * 启动方式:
 *   go run ./zy.go examples/websocket/index.php
 *
 * 测试方式（浏览器控制台）:
 *   const ws = new WebSocket("ws://127.0.0.1:9501/ws");
 *   ws.onmessage = (e) => console.log("recv:", e.data);
 *   ws.onopen = () => ws.send("hello origami");
 */

$server = new Server("0.0.0.0", port: 9501);

$server->get("/info", function ($request, $response) {
    $response->json([
        "name" => "Origami WebSocket Example",
        "ws" => "ws://127.0.0.1:9501/ws",
        "usage" => "connect then send text, server will echo",
    ]);
});

$server->get("/ws", function ($request, $response) {
    $conn = upgrade($request, $response);

    try {
        while (true) {
            $message = $conn->readText();
            if ($message == null || $message == "") {
                continue;
            }
            $conn->writeText("echo: " . $message);
        }
    } catch (Exception $e) {
        Log::info("websocket disconnected: " . $e->getMessage());
    } finally {
        $conn->close();
    }
});

Log::info("WebSocket 示例启动: http://127.0.0.1:9501/info");
Log::info("WebSocket 连接地址: ws://127.0.0.1:9501/ws");
$server->run();
