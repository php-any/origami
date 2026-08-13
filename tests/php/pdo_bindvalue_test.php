<?php

namespace tests\php;

/**
 * 验证 PDOStatement::bindValue + execute()（无参）可用于 INSERT。
 */

$path = __DIR__ . '/pdo_bindvalue_test.sqlite';
if (is_file($path)) {
    unlink($path);
}

$pdo = new \PDO('sqlite:' . $path);
$pdo->exec('CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)');
$stmt = $pdo->prepare('INSERT INTO users (name) VALUES (?)');
$stmt->bindValue(1, 'Alice');
$ok = $stmt->execute();
if ($ok !== true) {
    Log::fatal('pdo_bindvalue: execute 失败');
}
$id = $pdo->lastInsertId();
if ((int) $id < 1) {
    Log::fatal('pdo_bindvalue: lastInsertId 不对: ' . var_export($id, true));
}

@unlink($path);
Log::info('pdo_bindvalue 测试通过');
