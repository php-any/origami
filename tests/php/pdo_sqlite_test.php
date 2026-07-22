<?php

namespace tests\php;

/**
 * 验证 PDO sqlite 构造后状态正确落在实例上，可执行 DDL/DML。
 */

$path = __DIR__ . '/pdo_sqlite_test.sqlite';
if (is_file($path)) {
    unlink($path);
}

$pdo = new \PDO('sqlite:' . $path);
$pdo->exec('CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)');
$pdo->exec("INSERT INTO users (name) VALUES ('Alice')");
$stmt = $pdo->query('SELECT name FROM users');
$row = $stmt->fetch(\PDO::FETCH_ASSOC);

if (!is_array($row) || ($row['name'] ?? null) !== 'Alice') {
    Log::fatal('pdo_sqlite: fetch 结果不对: ' . var_export($row, true));
}

$id = $pdo->lastInsertId();
$idInt = (int) $id;
if ($idInt != 1) {
    Log::fatal('pdo_sqlite: lastInsertId 不对: ' . var_export($id, true));
}

@unlink($path);
Log::info('pdo_sqlite 测试通过');
