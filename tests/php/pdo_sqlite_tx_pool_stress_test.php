<?php

namespace tests\php;

/**
 * 压力测试：多次 BEGIN IMMEDIATE 事务，暴露 SQLite 连接池未限制 MaxOpenConns=1 时的锁冲突。
 */
$path = __DIR__ . '/pdo_sqlite_tx_pool_stress_test.sqlite';
@unlink($path);
@unlink($path . '-wal');
@unlink($path . '-shm');

$pdo = new \PDO('sqlite:' . $path, null, null, [
    \PDO::ATTR_ERRMODE => \PDO::ERRMODE_EXCEPTION,
]);
$pdo->exec('PRAGMA journal_mode=WAL');
$pdo->exec('PRAGMA busy_timeout=1000');
$pdo->exec('CREATE TABLE t(id INTEGER PRIMARY KEY AUTOINCREMENT, v TEXT)');

// Warm pool with extra idle connections by preparing on separate statements
$stmts = [];
for ($i = 0; $i < 8; $i++) {
    $stmts[] = $pdo->prepare('SELECT 1');
    $stmts[$i]->execute();
    $stmts[$i]->fetch();
}

$ok = 0;
try {
    for ($i = 0; $i < 20; $i++) {
        $pdo->exec('BEGIN IMMEDIATE');
        // force another connection touch before write
        $pdo->query('SELECT 1')->fetch();
        $s = $pdo->prepare('INSERT INTO t(v) VALUES(?)');
        $s->execute(['row' . $i]);
        $id = (int) $pdo->lastInsertId();
        $pdo->exec('COMMIT');
        if ($id < 1) {
            Log::fatal("stress lastInsertId=0 at i=$i");
        }
        $ok++;
    }
} catch (\Throwable $e) {
    Log::fatal("stress failed after ok=$ok: " . get_class($e) . ': ' . $e->getMessage());
}

$cnt = (int) $pdo->query('SELECT COUNT(*) FROM t')->fetchColumn();
if ($cnt !== 20) {
    Log::fatal("stress count=$cnt want 20");
}
@unlink($path);
@unlink($path . '-wal');
@unlink($path . '-shm');
Log::info('pdo_sqlite_tx_pool_stress 测试通过');
