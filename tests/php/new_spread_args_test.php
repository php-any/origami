<?php

namespace tests\php;

/**
 * 验证 new Class(...$args) 正确展开变参（修复 IncomingEntry::make 等内容嵌套问题）。
 */

class NewSpreadArgs_Sample
{
    public $content;
    public $uuid;

    public function __construct(array $content, $uuid = null)
    {
        $this->content = $content;
        $this->uuid = $uuid;
    }

    public static function make(...$arguments)
    {
        return new static(...$arguments);
    }
}

$direct = new NewSpreadArgs_Sample(['level' => 'info', 'message' => 'direct']);
if (($direct->content['message'] ?? null) !== 'direct') {
    Log::fatal('new Class(array) 失败');
}

$viaMake = NewSpreadArgs_Sample::make(['level' => 'info', 'message' => 'via make']);
if (($viaMake->content['message'] ?? null) !== 'via make') {
    Log::fatal('new static(...$arguments) 未展开: ' . var_export($viaMake->content, true));
}

$multi = NewSpreadArgs_Sample::make(['a' => 1], 'uid-1');
if (($multi->content['a'] ?? null) !== 1 || $multi->uuid !== 'uid-1') {
    Log::fatal('new static(...$arguments) 多参展开失败');
}

Log::info('new_spread_args 测试通过');
