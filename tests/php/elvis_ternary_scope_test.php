<?php

/**
 * 三元运算符 true 分支含 $var::{expr}()（Laravel HigherOrderCollectionProxy）。
 */
class ElvisTernaryScope_Target
{
    public static function ping(string $msg): string
    {
        return strtoupper($msg);
    }
}

$value = ElvisTernaryScope_Target::class;
$method = 'ping';
$parameters = ['hello'];

$result = is_string($value)
    ? $value::{$method}(...$parameters)
    : $value->{$method}(...$parameters);

if ($result !== 'HELLO') {
    \Log::fatal('ternary with :: in true branch failed: ' . var_export($result, true));
}
\Log::info('elvis_ternary_scope_test OK');
