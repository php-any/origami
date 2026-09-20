<?php

namespace tests\php;

/**
 * 复现 Filament tables index：array_filter(fn (Column $column)) 后 foreach $column->getName()。
 * 外层同名 $column 若被箭头函数捕获覆盖，会把整数传进 getName。
 */

class AFCol_GetName
{
    public function __construct(public string $name, public bool $sortable = true)
    {
    }

    public function isSortable(): bool
    {
        return $this->sortable;
    }

    public function isVisible(): bool
    {
        return true;
    }

    public function isToggledHidden(): bool
    {
        return false;
    }

    public function getName(): string
    {
        return $this->name;
    }
}

$column = 3;
$resultCount = 3;

$columns = [
    'name' => new AFCol_GetName('name'),
    'display_name' => new AFCol_GetName('display_name'),
    'permissions_count' => new AFCol_GetName('permissions_count', false),
];

$visible = array_filter(
    $columns,
    fn (AFCol_GetName $column): bool => $column->isVisible() && (! $column->isToggledHidden()),
);

$sortableColumns = array_filter(
    $visible,
    fn (AFCol_GetName $column): bool => $column->isSortable(),
);

if (count($sortableColumns) !== 2) {
    \Log::fatal('sortable 数量不对: '.count($sortableColumns));
}

$names = [];
foreach ($sortableColumns as $column) {
    if (!is_object($column)) {
        \Log::fatal('foreach $column 不是对象: '.var_export($column, true));
    }
    $names[] = $column->getName();
}

if ($names !== ['name', 'display_name']) {
    \Log::fatal('getName 结果不对: '.var_export($names, true));
}

// 外层 leftover 不应再是 3
if (!is_object($column) || $column->getName() !== 'display_name') {
    \Log::fatal('foreach 结束后 $column leftover: '.var_export($column, true));
}

\Log::info('array_filter typed arrow foreach 测试通过');
