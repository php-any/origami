<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Database\Capsule\Manager as Capsule;
use Illuminate\Support\Arr;
use Illuminate\Support\Fluent;

/**
 * 仅验证 Capsule 连接与简单查询；schema 迁移依赖更完整的 PDO/异常语义。
 */
$capsule = new Capsule();

$config = new Fluent([
    'database.fetch' => PDO::FETCH_OBJ,
    'database.default' => 'default',
    'database.connections' => [
        'default' => [
            'driver' => 'sqlite',
            'database' => ':memory:',
            'prefix' => '',
        ],
    ],
]);
$capsule->getContainer()->instance('config', $config);

$capsule->setAsGlobal();
$capsule->bootEloquent();

$conn = $capsule->getConnection();
$rows = $conn->select('select 1 as n');
$row = $rows[0] ?? null;
$n = Arr::get($row, 'n');
if ($n != 1) {
    echo "FAIL: select, got ";
    var_export($row);
    echo "\n";
    exit(1);
}

echo "PASS\n";
