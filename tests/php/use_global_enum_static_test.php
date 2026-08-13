<?php

/**
 * use 导入全局 enum 后，命名空间内访问 Enum::Case 不得被错误加上当前命名空间前缀。
 * 回归：Laravel 13 SortDirection（use SortDirection; … SortDirection::Descending）。
 */

include __DIR__ . '/use_global_enum_static_direction.php';

namespace UseGlobalEnumStatic_App;

use UseGlobalEnumStatic_Direction;

function directionLabel($d)
{
    return match (true) {
        $d instanceof UseGlobalEnumStatic_Direction => match ($d) {
            UseGlobalEnumStatic_Direction::Ascending => 'asc',
            UseGlobalEnumStatic_Direction::Descending => 'desc',
        },
        default => 'unknown',
    };
}

$desc = UseGlobalEnumStatic_Direction::Descending;
$label = directionLabel($desc);
if ($label !== 'desc') {
    \Log::fatal('use_global_enum_static: Descending 应为 desc, got ' . var_export($label, true));
}

$asc = UseGlobalEnumStatic_Direction::Ascending;
$label = directionLabel($asc);
if ($label !== 'asc') {
    \Log::fatal('use_global_enum_static: Ascending 应为 asc, got ' . var_export($label, true));
}

\Log::info('use_global_enum_static 测试通过');
