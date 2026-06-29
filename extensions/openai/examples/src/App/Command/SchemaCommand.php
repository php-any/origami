<?php

namespace App\Command;

use Cli\Annotation\Command;

include __DIR__ . "/load_env.php";

#[Command(name: "schema", description: "JSON Schema — 严格模式 + json_object 降级")]
class SchemaCommand
{
    public function execute(): void
    {
        $env = loadEnv();

        Log::info("=== JSON Schema ===");
        Log::info("MODEL: {$env['model']}");

        $c = new \OpenAI\Client($env['key'], $env['base']);
        $m = $env['model'];

        // 1. json_schema 严格模式 (仅 OpenAI 支持)
        Log::info("--- 1. json_schema ---");
        $schema = [
            "type" => "object",
            "properties" => [
                "name" => ["type" => "string"],
                "age"  => ["type" => "integer"],
                "city" => ["type" => "string"],
            ],
            "required" => ["name", "age", "city"],
            "additionalProperties" => false,
        ];
        try {
            $r = $c->chat($m, [
                new \OpenAI\SystemMessage("根据用户输入提取信息，严格按 JSON Schema 输出。"),
                new \OpenAI\UserMessage("小李今年25岁，住在上海。"),
            ], [
                \OpenAI\ChatOption::RESPONSE_FORMAT => [
                    "type"   => "json_schema",
                    "name"   => "person",
                    "strict" => true,
                    "schema" => $schema,
                ],
                \OpenAI\ChatOption::TEMPERATURE => 0,
            ]);
            $d = json_decode($r->content);
            Log::info("✓ name: " . $d->name . ", age: " . $d->age . ", city: " . $d->city);
        } catch (\Exception $e) {
            Log::info("⚠ json_schema 不支持: " . $e->getMessage());
        }

        // 2. 降级: json_object + 示例驱动
        Log::info("--- 2. 降级 json_object ---");
        $sys = "Extract person info from user input. Output ONLY valid JSON.\n\nEXAMPLE INPUT: 小李今年25岁，住在上海。\nEXAMPLE JSON OUTPUT:\n{\"name\": \"小李\", \"age\": 25, \"city\": \"上海\"}";
        $r = $c->chat($m, [
            new \OpenAI\SystemMessage($sys),
            new \OpenAI\UserMessage("小王30岁，住在北京。"),
        ], [\OpenAI\ChatOption::RESPONSE_FORMAT => "json_object", \OpenAI\ChatOption::TEMPERATURE => 0]);
        $d = json_decode($r->content);
        Log::info("✓ name: " . $d->name . ", age: " . $d->age . ", city: " . $d->city);

        Log::info("=== 完成 ===");
    }
}
