<?php

namespace App\Command;

use Cli\Annotation\Command;

include __DIR__ . "/load_env.php";

#[Command(name: "chat", description: "Chat Completions — 消息类、多轮对话、options")]
class ChatCommand
{
    public function execute(): void
    {
        $env = loadEnv();

        Log::info("=== Chat Completions ===");
        Log::info("MODEL: {$env['model']}");

        $c = new \OpenAI\Client($env['key'], $env['base']);
        $m = $env['model'];

        // 1. UserMessage
        Log::info("--- 1. UserMessage ---");
        $r = $c->chat($m, [new \OpenAI\UserMessage("1+1等于几？一句话回答")]);
        Log::info("回复: " . $r->content . " (tokens: " . $r->usage->totalTokens . ")");

        // 2. SystemMessage + ChatOption
        Log::info("--- 2. SystemMessage + ChatOption ---");
        $r = $c->chat($m, [
            new \OpenAI\SystemMessage("你只输出 JSON，不要有其他文字。"),
            new \OpenAI\UserMessage("返回 {\"answer\": 42}"),
        ], [
            \OpenAI\ChatOption::TEMPERATURE => 0,
            \OpenAI\ChatOption::SEED => 42,
        ]);
        Log::info("回复: " . $r->content);

        // 3. 多轮对话
        Log::info("--- 3. 多轮对话 ---");
        $r = $c->chat($m, [
            new \OpenAI\UserMessage("我叫小明"),
            new \OpenAI\AssistantMessage("你好小明！"),
            new \OpenAI\UserMessage("我叫什么？"),
        ]);
        Log::info("回复: " . $r->content);

        // 4. maxTokens
        Log::info("--- 4. MAX_TOKENS=20 ---");
        $r = $c->chat($m, [new \OpenAI\UserMessage("写一篇长文章")], [
            \OpenAI\ChatOption::MAX_TOKENS => 20,
        ]);
        Log::info("finishReason: " . $r->finishReason);

        // 5. DeveloperMessage (仅 OpenAI)
        Log::info("--- 5. DeveloperMessage ---");
        try {
            $r = $c->chat($m, [
                new \OpenAI\DeveloperMessage("你是命令行工具，只输出结果。"),
                new \OpenAI\UserMessage("计算 3*4"),
            ]);
            Log::info("回复: " . $r->content);
        } catch (\Exception $e) {
            Log::info("不支持: " . $e->getMessage());
        }

        // 6. responseFormat
        Log::info("--- 6. RESPONSE_FORMAT ---");
        $r = $c->chat($m, [
            new \OpenAI\SystemMessage("只输出 JSON。"),
            new \OpenAI\UserMessage("返回 {\"status\": \"ok\"}"),
        ], [
            \OpenAI\ChatOption::RESPONSE_FORMAT => "json_object",
            \OpenAI\ChatOption::TEMPERATURE => 0,
        ]);
        Log::info("回复: " . $r->content);

        Log::info("=== 完成 ===");
    }
}
