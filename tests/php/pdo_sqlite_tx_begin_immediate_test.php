<?php

namespace tests\php;

/**
 * 复现 bbs1org 发帖事务路径：exec(BEGIN IMMEDIATE) + prepare/execute INSERT + lastInsertId + COMMIT。
 */

$path = __DIR__ . '/pdo_sqlite_tx_begin_immediate_test.sqlite';
@unlink($path);
@unlink($path . '-wal');
@unlink($path . '-shm');

$pdo = new \PDO('sqlite:' . $path, null, null, [
    \PDO::ATTR_ERRMODE => \PDO::ERRMODE_EXCEPTION,
    \PDO::ATTR_DEFAULT_FETCH_MODE => \PDO::FETCH_ASSOC,
    \PDO::ATTR_EMULATE_PREPARES => false,
]);
foreach (['PRAGMA journal_mode=WAL', 'PRAGMA busy_timeout=5000', 'PRAGMA foreign_keys=ON'] as $sql) {
    $pdo->exec($sql);
}
$pdo->exec('CREATE TABLE app_topics(id INTEGER PRIMARY KEY AUTOINCREMENT,forum_id INTEGER,user_id INTEGER,title TEXT,body TEXT,created_at INTEGER,last_reply_at INTEGER)');

function pdo_sqlite_tx_q(\PDO $pdo, string $sql, array $p = []): \PDOStatement
{
    $s = $pdo->prepare($sql);
    $s->execute($p);
    return $s;
}

try {
    $pdo->exec('BEGIN IMMEDIATE');
    pdo_sqlite_tx_q($pdo, 'INSERT INTO app_topics(forum_id,user_id,title,body,created_at,last_reply_at) VALUES(?,?,?,?,?,?)', [1, 1, 't', 'b', 123, 123]);
    $tid = (int) $pdo->lastInsertId();
    $pdo->exec('COMMIT');
    $cnt = (int) $pdo->query('SELECT COUNT(*) FROM app_topics')->fetchColumn();
    if ($tid < 1) {
        Log::fatal('pdo_sqlite_tx: lastInsertId 不对: ' . $tid);
    }
    if ($cnt < 1) {
        Log::fatal('pdo_sqlite_tx: 提交后行数为 0');
    }
    Log::info("pdo_sqlite_tx BEGIN IMMEDIATE 路径通过 tid=$tid count=$cnt");
} catch (\Throwable $e) {
    try {
        $pdo->exec('ROLLBACK');
    } catch (\Throwable $e2) {
    }
    Log::fatal('pdo_sqlite_tx 异常: ' . get_class($e) . ': ' . $e->getMessage());
}

@unlink($path);
@unlink($path . '-wal');
@unlink($path . '-shm');
