<?php

namespace App\Command;

use Cli\Annotation\Command;

include __DIR__ . "/load_env.php";

#[Command(name: "error", description: "错误处理 — 401/400、类型错误、文件不存在")]
class ErrorCommand
{
    public function execute(): void
    {
        $env = loadEnv();

        Log::info("=== 错误处理 ===");

        // 1. 无效 API Key
        Log::info("--- 1. 无效 API Key ---");
        $c = new \OpenAI\Client("invalid-key-xxxxx", $env['base']);
        try {
            $c->chat("gpt-4o", [new \OpenAI\UserMessage("Hi")]);
            Log::error("应该抛出异常！");
        } catch (\Exception $e) {
            Log::info("✓ 401");
        }

        // 2. 无效模型名
        Log::info("--- 2. 无效模型名 ---");
        $c = new \OpenAI\Client($env['key'], $env['base']);
        try {
            $c->chat("this-model-does-not-exist", [new \OpenAI\UserMessage("Hi")]);
            Log::error("应该抛出异常！");
        } catch (\Exception $e) {
            Log::info("✓ 400");
        }

        // 3. model 类型错误
        Log::info("--- 3. model 类型错误 ---");
        try {
            $c->chat(12345, [new \OpenAI\UserMessage("Hi")]);
            Log::error("应该抛出异常！");
        } catch (\Exception $e) {
            Log::info("✓ 类型错误");
        }

        // 4. 空消息数组
        Log::info("--- 4. 空消息数组 ---");
        try {
            $c->chat("gpt-4o", []);
            Log::info("空消息已发送（API 端校验）");
        } catch (\Exception $e) {
            Log::info("✓ " . $e->getMessage());
        }

        // 5. transcription 文件不存在
        Log::info("--- 5. transcription 文件不存在 ---");
        try {
            $c->transcription("whisper-1", "/nonexistent/audio.mp3");
            Log::error("应该抛出异常！");
        } catch (\Exception $e) {
            Log::info("✓ 文件错误");
        }

        Log::info("=== 完成 ===");
    }
}
