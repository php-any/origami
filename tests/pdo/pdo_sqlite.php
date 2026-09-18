<?php

namespace tests\pdo;

/**
 * pdo：ATTR_ERRMODE 与 sqlite 内存库（驱动可用时）。
 */

$mode = \PDO::ATTR_ERRMODE;
if ($mode === null || $mode === false || $mode === '') {
    Log::fatal('PDO::ATTR_ERRMODE 缺失');
}
if (method_exists('PDO', 'getAvailableDrivers')) {
    $drivers = \PDO::getAvailableDrivers();
    if (is_array($drivers) && in_array('sqlite', $drivers, true)) {
        $pdo = new \PDO('sqlite::memory:');
        $pdo->exec('CREATE TABLE t (id INTEGER)');
        $n = $pdo->exec('INSERT INTO t VALUES (1)');
        if ($n === false) {
            Log::fatal('sqlite 内存插入失败');
        }
    }
}

Log::info('pdo sqlite/attr 测试通过');
