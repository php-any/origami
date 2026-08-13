<?php

namespace tests\php;

/**
 * 验证 preg_match 的 offset 与 /A 锚定修饰符（Symfony Dotenv 解析 .env 依赖此行为）。
 */

$subject = "# GitHub OAuth（登录）\nGITHUB_CLIENT_ID=\"secret\"\n";
$matches = [];

// 无 offset 时会在注释行误匹配 "GitHub"
$count = preg_match('/([A-Z][A-Z0-9_]*+)/A', $subject, $matches);
if ($count !== 0) {
    Log::fatal('无 offset 时 /A 不应在注释后匹配变量名, got=' . json_encode($matches));
}

// 从第二行起应匹配 GITHUB_CLIENT_ID
$offset = strpos($subject, 'GITHUB');
$count = preg_match('/([A-Z][A-Z0-9_]*+)/A', $subject, $matches, 0, $offset);
if ($count !== 1 || $matches[1] !== 'GITHUB_CLIENT_ID') {
    Log::fatal('offset+/A 匹配失败: count=' . $count . ' matches=' . json_encode($matches));
}

// skipEmptyLines 风格：跳过注释行后解析变量名
$cursor = 0;
preg_match('/(?:\s*(?:#[^\n]*)?)++/A', $subject, $skip, 0, $cursor);
$cursor += strlen($skip[0]);

$count = preg_match('/(export[ \t]++)?([A-Z][A-Z0-9_]*+)/A', $subject, $matches, 0, $cursor);
if ($count !== 1 || $matches[2] !== 'GITHUB_CLIENT_ID') {
    Log::fatal('dotenv lexVarname 风格匹配失败: count=' . $count . ' matches=' . json_encode($matches));
}

Log::info('preg_match offset+/A 测试通过');
