<?php

namespace Spring\Service;

use Container\Singleton;

/**
 * WebSocket 连接池：维护活跃连接并支持广播
 */
#[Singleton]
class WebSocketHub {

    /** @var array<int, mixed> */
    private array $connections = [];

    public function add($conn): void {
        $this->connections[] = $conn;
        Log::info("WebSocket 连接建立，当前在线: " . $this->count());
    }

    public function remove($conn): void {
        foreach ($this->connections as $index => $item) {
            if ($item === $conn) {
                unset($this->connections[$index]);
                break;
            }
        }
        $this->connections = array_values($this->connections);
        Log::info("WebSocket 连接关闭，当前在线: " . $this->count());
    }

    public function count(): int {
        return count($this->connections);
    }

    public function broadcast(string $message, $except = null): void {
        $dead = [];
        foreach ($this->connections as $index => $conn) {
            if ($except !== null && $conn === $except) {
                continue;
            }
            try {
                $conn->writeText($message);
            } catch (\Throwable $e) {
                $dead[] = $index;
            }
        }
        foreach ($dead as $index) {
            unset($this->connections[$index]);
        }
        if (count($dead) > 0) {
            $this->connections = array_values($this->connections);
        }
    }
}
