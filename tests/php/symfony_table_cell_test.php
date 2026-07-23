<?php

namespace tests\php;

/**
 * 复现 Symfony Console Table：纯字符串单元格却报 object（fillNextRows / is_scalar）。
 */

require dirname(__DIR__, 2) . '/examples/laravel/vendor/autoload.php';

use Symfony\Component\Console\Helper\Table;
use Symfony\Component\Console\Helper\TableSeparator;
use Symfony\Component\Console\Output\BufferedOutput;

$rows = [
    ['1', 'Alice', 'a@b.com'],
];

// 1) 直接 is_scalar
foreach ($rows[0] as $c) {
    if (!is_scalar($c)) {
        Log::fatal('plain cell should be scalar, got ' . get_debug_type($c));
    }
}

// 2) 模拟 array_merge(headers, [divider], rows) 后的 fillNextRows 检查
$divider = new TableSeparator();
$merged = array_merge([['ID', 'Name', 'Email']], [$divider], $rows);

Log::info('merged count=' . count($merged));
for ($i = 0; $i < count($merged); $i++) {
    $row = $merged[$i];
    $isSep = $row instanceof TableSeparator;
    Log::info("row $i type=" . get_debug_type($row) . ' instanceof_sep=' . ($isSep ? '1' : '0'));

    if ($isSep) {
        // PHP 中 foreach TableSeparator 会怎样？
        $n = 0;
        foreach ($row as $col => $cell) {
            $n++;
            Log::info("  sep foreach col=$col type=" . get_debug_type($cell) . ' scalar=' . (is_scalar($cell) ? '1' : '0'));
            if ($cell !== null && !($cell instanceof \Symfony\Component\Console\Helper\TableCell) && !is_scalar($cell) && !($cell instanceof \Stringable)) {
                Log::info('  FAIL-like cell at sep foreach');
            }
        }
        Log::info("  sep foreach count=$n");
        continue;
    }

    foreach ($row as $col => $cell) {
        $ok = is_scalar($cell) || $cell === null;
        Log::info("  cell col=$col type=" . get_debug_type($cell) . ' scalar=' . (is_scalar($cell) ? '1' : '0'));
        if (!$ok && !($cell instanceof \Symfony\Component\Console\Helper\TableCell) && !($cell instanceof \Stringable)) {
            Log::fatal("non-scalar cell at row $i col $col: " . get_debug_type($cell));
        }
    }
}

// 3) 真实 Table::render
$out = new BufferedOutput();
$t = new Table($out);
$t->setHeaders(['ID', 'Name', 'Email']);
$t->setRows($rows);
try {
    $t->render();
    $text = $out->fetch();
    if ($text === '') {
        Log::fatal('Table::render produced empty output');
    }
    Log::info("Table::render ok, len=" . strlen($text));
} catch (\Throwable $e) {
    Log::fatal('Table::render threw: ' . $e->getMessage());
}

Log::info('symfony_table_cell_test 测试通过');
