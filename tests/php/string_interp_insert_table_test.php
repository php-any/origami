<?php

namespace tests\php;

/**
 * Laravel compileInsert 使用 "insert into $table ($columns) values $parameters"。
 * MySQL wrap 后 $table 含反引号，插值不能把 `users` 变成 `user)`。
 */

$table = '`users`';
$columns = '`name`, `email`, `password`, `updated_at`, `created_at`';
$parameters = '(?, ?, ?, ?, ?)';
$sql = "insert into $table ($columns) values $parameters";
$expected = 'insert into `users` (`name`, `email`, `password`, `updated_at`, `created_at`) values (?, ?, ?, ?, ?)';
if ($sql !== $expected) {
    Log::fatal("compileInsert 风格插值失败: $sql");
}

$parts = preg_split('/(.)(?=[A-Z])/u', 'User', -1, PREG_SPLIT_DELIM_CAPTURE);
if ($parts === false) {
    Log::fatal('preg_split 前瞻断言应能编译，不能返回 false');
}
$lastWord = array_pop($parts);
$joined = implode('', $parts).$lastWord;
if ($joined !== 'User') {
    Log::fatal('pluralStudly 用的 preg_split 拆坏了 User: ' . json_encode(['parts' => $parts, 'last' => $lastWord, 'joined' => $joined]));
}

Log::info('string_interp_insert_table 测试通过');
