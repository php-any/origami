<?php

namespace tests\php;

/**
 * PHP 8：未定义方法的命名实参必须进 __call 的 $arguments 关联键。
 * Filament DatabaseNotifications：Relation::simplePaginate(50, pageName: '...')
 * 经 __call 转发；丢掉键会把 pageName 当成 $columns，SQL 选出不存在的列。
 */

class MagicCallNamed_Rel
{
    public function __call($method, $parameters)
    {
        return $parameters;
    }
}

$p = (new MagicCallNamed_Rel())->simplePaginate(50, pageName: 'database-notifications-page');
if (!is_array($p) || ($p[0] ?? null) !== 50) {
    \Log::fatal('位置实参 50 丢失: ' . var_export($p, true));
}
if (!array_key_exists('pageName', $p) || $p['pageName'] !== 'database-notifications-page') {
    \Log::fatal('pageName 必须为关联键, 实际: ' . var_export($p, true));
}
if (array_key_exists(1, $p)) {
    \Log::fatal('pageName 不应落到位置 1: ' . var_export($p, true));
}

class MagicCallNamed_Target
{
    public function simplePaginate($perPage = 15, $columns = ['*'], $pageName = 'page', $page = null)
    {
        return compact('perPage', 'columns', 'pageName', 'page');
    }
}

class MagicCallNamed_Fwd
{
    public function __call($method, $parameters)
    {
        return (new MagicCallNamed_Target())->$method(...$parameters);
    }
}

$r = (new MagicCallNamed_Fwd())->simplePaginate(50, pageName: 'database-notifications-page');
if (($r['perPage'] ?? null) !== 50) {
    \Log::fatal('转发后 perPage 错误: ' . var_export($r, true));
}
if (($r['columns'] ?? null) !== ['*']) {
    \Log::fatal('转发后 columns 应为默认 [*], 实际: ' . var_export($r, true));
}
if (($r['pageName'] ?? null) !== 'database-notifications-page') {
    \Log::fatal('转发后 pageName 错误: ' . var_export($r, true));
}

\Log::info('__call named args 测试通过');
