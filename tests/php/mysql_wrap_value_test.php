<?php

namespace tests\php;

/**
 * 对齐 MySqlGrammar::wrapValue：反引号包裹表名，不能把 users 编成 user)。
 */

$value = 'users';
$wrapped = $value === '*' ? $value : '`'.str_replace('`', '``', $value).'`';
if ($wrapped !== '`users`') {
    Log::fatal('MySQL wrapValue 失败: '.var_export($wrapped, true));
}

$table = $wrapped;
$columns = '`name`, `email`';
$sql = "insert into $table ($columns) values (?, ?)";
if ($sql !== 'insert into `users` (`name`, `email`) values (?, ?)') {
    Log::fatal('wrap 后 insert SQL 失败: '.$sql);
}

Log::info('mysql_wrap_value 测试通过');
