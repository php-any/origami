<?php

namespace App\Command;

use Cli\Annotation\Command;

include __DIR__ . "/load_env.php";

#[Command(name: "json", description: "JSON 输出 — json_object 模式、示例驱动提取")]
class JsonCommand
{
    public function execute(): void
    {
        $env = loadEnv();

        Log::info("=== JSON 输出 ===");
        Log::info("MODEL: {$env['model']}");

        $c = new \OpenAI\Client($env['key'], $env['base']);
        $m = $env['model'];
        $opt = [\OpenAI\ChatOption::RESPONSE_FORMAT => "json_object", \OpenAI\ChatOption::TEMPERATURE => 0];

        // 1. 基本 json_object
        Log::info("--- 1. json_object 基本 ---");
        $r = $c->chat($m, [
            new \OpenAI\SystemMessage("你只输出 JSON，不要有其他文字。"),
            new \OpenAI\UserMessage("返回 {\"name\": \"John\", \"age\": 30}"),
        ], $opt);
        $d = json_decode($r->content);
        Log::info("✓ name=" . $d->name . ", age=" . $d->age);

        // 2. QA 提取 — 对应 DeepSeek Python 示例
        Log::info("--- 2. QA 提取 ---");
        $sys = "The user will provide some exam text. Please parse the \"question\" and \"answer\" and output them in JSON format.\n\nEXAMPLE INPUT:\nWhich is the highest mountain in the world? Mount Everest.\n\nEXAMPLE JSON OUTPUT:\n{\"question\": \"Which is the highest mountain in the world?\", \"answer\": \"Mount Everest\"}";
        $r = $c->chat($m, [
            new \OpenAI\SystemMessage($sys),
            new \OpenAI\UserMessage("Which is the longest river in the world? The Nile River."),
        ], $opt);
        $d = json_decode($r->content);
        Log::info("✓ question: " . $d->question);
        Log::info("✓ answer: " . $d->answer);

        // 3. 多字段提取
        Log::info("--- 3. 多字段提取 ---");
        $sys = "Extract structured data from the user input. Output ONLY valid JSON.\n\nEXAMPLE INPUT:\nMy name is Alice, I am 28 years old, I know Python and JavaScript.\n\nEXAMPLE JSON OUTPUT:\n{\"name\": \"Alice\", \"age\": 28, \"skills\": [\"Python\", \"JavaScript\"]}";
        $r = $c->chat($m, [
            new \OpenAI\SystemMessage($sys),
            new \OpenAI\UserMessage("I am Bob, 35 years old, I know Go, Rust and PHP."),
        ], $opt);
        $d = json_decode($r->content);
        Log::info("✓ name: " . $d->name . ", age: " . $d->age . ", skills: " . count($d->skills));

        // 4. 产品信息提取
        Log::info("--- 4. 产品信息 ---");
        $r = $c->chat($m, [
            new \OpenAI\SystemMessage("从描述中提取产品信息，只输出 JSON。"),
            new \OpenAI\UserMessage("iPhone 15 Pro, \$999, in stock, category: phones, tags: apple,5G"),
        ], $opt);
        Log::info("✓ " . $r->content);

        Log::info("=== 完成 ===");
    }
}
