<?php

namespace Spring\Controller;

use Net\Annotation\Controller;
use Net\Annotation\GetMapping;
use Net\Annotation\Route;
use Spring\Service\WebSocketHub;
use function Net\Websocket\upgrade;

/**
 * WebSocket 聊天室端点
 *
 * 客户端连接 ws://host:port/ws/chat，发送 JSON 文本消息：
 * {"type":"message","from":"昵称","text":"你好"}
 */
#[Controller]
#[Route(prefix: "/ws")]
class WebSocketController {

    public function __construct(
        private WebSocketHub $hub,
    ) {}

    #[GetMapping(path: "/chat")]
    public function chat($request, $response):void {
        $conn = upgrade($request, $response, true);
        $this->hub->add($conn);

        $joinMsg = json_encode([
            'type' => 'system',
            'text' => '新用户加入聊天室',
            'online' => $this->hub->count(),
            'time' => time(),
        ]);

        try {
            $conn->writeText(json_encode([
                'type' => 'system',
                'text' => '欢迎加入 Spring WebSocket 聊天室',
                'online' => $this->hub->count(),
                'time' => time(),
            ]));
            $this->hub->broadcast($joinMsg, $conn);

            while (true) {
                $raw = $conn->readText();
                if ($raw === '') {
                    continue;
                }

                $payload = json_decode($raw, true);
                $type = 'message';
                $from = '匿名';
                $text = $raw;

                if (is_array($payload)) {
                    $type = $payload['type'] ?? 'message';
                    $from = $payload['from'] ?? '匿名';
                    $text = $payload['text'] ?? '';
                }

                $this->hub->broadcast(json_encode([
                    'type' => $type,
                    'from' => $from,
                    'text' => $text,
                    'time' => time(),
                    'online' => $this->hub->count(),
                ]));
            }
        } catch (\Throwable $e) {
            \Log::info("WebSocket 会话结束: " . $e->getMessage());
        } finally {
            $this->hub->remove($conn);
            try {
                $conn->close();
            } catch (\Throwable $e) {
                // 连接可能已关闭
            }
            $this->hub->broadcast(json_encode([
                'type' => 'system',
                'text' => '用户离开聊天室',
                'online' => $this->hub->count(),
                'time' => time(),
            ]));
        }
    }
}
