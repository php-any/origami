<?php

namespace tests\php;

/**
 * preg_replace 失败/无匹配时不应变成布尔 false 再被 str_replace 变成 "false"。
 */

$s = 'admin.login';
$out = preg_replace('/⚡[\x{FE0E}\x{FE0F}]?/u', '', $s);
if ($out !== 'admin.login') {
    Log::fatal('preg_replace 无匹配应返回原串, got=' . var_export($out, true) . ' type=' . gettype($out));
}

$out2 = str_replace('/', '.', $out);
if ($out2 !== 'admin.login') {
    Log::fatal('str_replace after preg_replace 失败: ' . var_export($out2, true));
}

// is_subclass_of 对普通字符串应为 false
if (is_subclass_of('admin.login', \Exception::class)) {
    Log::fatal('is_subclass_of(admin.login, Exception) 不应为 true');
}

Log::info('preg_replace/normalize 相关测试通过: out=' . $out);
