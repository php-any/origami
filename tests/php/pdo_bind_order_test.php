<?php

namespace tests\php;

/**
 * bindValue 位置参数后 execute() 必须按 1..n 顺序绑定。
 */

$fail = 0;
for ($i = 0; $i < 50; $i++) {
    $pdo = new \PDO('sqlite::memory:');
    $pdo->exec('CREATE TABLE t (a TEXT, b TEXT)');
    $stmt = $pdo->prepare('INSERT INTO t (a, b) VALUES (?, ?)');
    $bindings = ['Alice', 'alice@example.com'];
    foreach ($bindings as $key => $value) {
        $stmt->bindValue($key + 1, $value);
    }
    $stmt->execute();
    $row = $pdo->query('SELECT a, b FROM t')->fetch(\PDO::FETCH_ASSOC);
    if (($row['a'] ?? null) !== 'Alice' || ($row['b'] ?? null) !== 'alice@example.com') {
        $fail++;
        if ($fail === 1) {
            Log::info('first bad: ' . json_encode($row));
        }
    }
}
if ($fail > 0) {
    Log::fatal("pdo_bind_order: $fail/50 绑定错位");
}
Log::info('pdo_bind_order 测试通过');
