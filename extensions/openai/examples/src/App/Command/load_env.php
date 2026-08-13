<?php

/**
 * 加载 .env 配置
 * 每个命令类通过 include 此文件获取 loadEnv() 函数
 */
function loadEnv(): array
{
    $dir = __DIR__;
    $content = file_get_contents(OS::path($dir, "..", "..", "..", ".env"));
    $lines = explode("\n", $content);
    $key = "";
    $base = "";
    $model = "";
    foreach ($lines as $line) {
        $line = $line->trim();
        if ($line === "" || $line->startsWith("#")) {
            continue;
        }
        $p = $line->indexOf("=");
        if ($p > 0) {
            $k = $line->substring(0, $p)->trim();
            $v = $line->substring($p + 1)->trim()->replace("\"", "")->replace("'", "");
            if ($k === "KEY") {
                $key = $v;
            }
            if ($k === "BASE_URL") {
                $base = $v;
            }
            if ($k === "MODEL") {
                $model = $v;
            }
        }
    }
    return ["key" => $key, "base" => $base, "model" => $model];
}
