<?php

namespace tests\php;

/** 验证 prepare/execute 数组参数与 lastInsertId。 */
$path = __DIR__ . '/pdo_sqlite_insert_params_test.sqlite';
@unlink($path);
$pdo = new \PDO('sqlite:' . $path);
$pdo->setAttribute(\PDO::ATTR_ERRMODE, \PDO::ERRMODE_EXCEPTION);
$pdo->exec('CREATE TABLE t(id INTEGER PRIMARY KEY AUTOINCREMENT, a INT, b TEXT)');
$s = $pdo->prepare('INSERT INTO t(a,b) VALUES(?,?)');
$s->execute([1, 'hello']);
$id = (int) $pdo->lastInsertId();
$row = $pdo->query('SELECT a,b FROM t')->fetch(\PDO::FETCH_ASSOC);
if ($id !== 1) {
    Log::fatal('insert_params lastInsertId=' . var_export($id, true));
}
if (!is_array($row) || (string)($row['a'] ?? '') !== '1' || (string)($row['b'] ?? '') !== 'hello') {
    Log::fatal('insert_params row=' . var_export($row, true));
}
@unlink($path);
Log::info('pdo_sqlite_insert_params 测试通过');
