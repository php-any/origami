<?php

namespace tests\php;

/**
 * PDO 从共享连接池借连接：同 DSN 多 PDO 可并存；close 归还连接后不可再查。
 */

$path = __DIR__ . '/pdo_pool_borrow_return_test.sqlite';
@unlink($path);
@unlink($path . '-wal');
@unlink($path . '-shm');
$dsn = 'sqlite:' . $path;

$a = new \PDO($dsn);
$a->setAttribute(\PDO::ATTR_ERRMODE, \PDO::ERRMODE_EXCEPTION);
$a->exec('CREATE TABLE t(id INTEGER PRIMARY KEY AUTOINCREMENT, v TEXT)');
$a->exec("INSERT INTO t(v) VALUES('a')");

// 第二个 PDO 从同一共享池再借一条连接（多线程场景）
$b = new \PDO($dsn);
$b->setAttribute(\PDO::ATTR_ERRMODE, \PDO::ERRMODE_EXCEPTION);
$row = $b->query('SELECT v FROM t')->fetch(\PDO::FETCH_ASSOC);
if (!is_array($row) || (string)($row['v'] ?? '') !== 'a') {
    Log::fatal('pool 第二 PDO 读不到第一 PDO 写入的数据: ' . var_export($row, true));
}

$b->exec("INSERT INTO t(v) VALUES('b')");
$cnt = (int) $a->query('SELECT COUNT(*) FROM t')->fetchColumn();
if ($cnt !== 2) {
    Log::fatal("pool 跨 PDO 可见性 count=$cnt");
}

// beginTransaction 落在借出的那条连接上
$a->beginTransaction();
$a->exec("INSERT INTO t(v) VALUES('tx')");
if (!$a->inTransaction()) {
    Log::fatal('inTransaction 应为 true');
}
$a->commit();
if ($a->inTransaction()) {
    Log::fatal('commit 后 inTransaction 应为 false');
}

$a->close();
try {
    $a->query('SELECT 1');
    Log::fatal('close 后仍能 query');
} catch (\Throwable $e) {
    // expected
}

// b 仍可用（连接已从池中借出，未受 a.close 影响）
$cnt2 = (int) $b->query('SELECT COUNT(*) FROM t')->fetchColumn();
if ($cnt2 !== 3) {
    Log::fatal("close 后另一 PDO count=$cnt2");
}
$b->close();

@unlink($path);
@unlink($path . '-wal');
@unlink($path . '-shm');
Log::info('pdo_pool_borrow_return 测试通过');
