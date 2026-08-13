<?php

/**
 * 应用配置
 *
 * 从项目根目录的 .env 读取 OpenAI 连接信息，返回配置数组。
 * 用法：$config = require __DIR__ . "/config/app.php";
 */

$defaults = [
    "key" => "",
    "base_url" => "https://api.openai.com/v1",
    "model" => "gpt-4o-mini",
];

$envFile = OS::path(__DIR__, "..", ".env");
if (is_file($envFile)) {
    $lines = explode("\n", file_get_contents($envFile));
    foreach ($lines as $line) {
        $line = $line->trim();
        if ($line === "" || $line->startsWith("#")) {
            continue;
        }
        $pos = $line->indexOf("=");
        if ($pos <= 0) {
            continue;
        }
        $name = $line->substring(0, $pos)->trim();
        $value = $line->substring($pos + 1)->trim()->replace("\"", "")->replace("'", "");

        if ($name === "KEY") {
            $defaults["key"] = $value;
        } elseif ($name === "BASE_URL") {
            $defaults["base_url"] = $value;
        } elseif ($name === "MODEL") {
            $defaults["model"] = $value;
        }
    }
}

return [
    "openai" => $defaults,
];
