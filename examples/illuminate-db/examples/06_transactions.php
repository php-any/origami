<?php
/**
 * 示例 6：数据库事务
 *
 * 演示使用 Illuminate Database 进行事务处理，
 * 包含事务提交与回滚两种场景。
 *
 * 运行方式: ./illuminate-db examples/06_transactions.php
 */

require __DIR__ . '/../vendor/autoload.php';
require __DIR__ . '/../bootstrap.php';

use Illuminate\Database\Capsule\Manager as Capsule;

// 初始化 Capsule
bootstrap_capsule();

// 创建账户表
Capsule::schema()->create('accounts', function ($table) {
    $table->increments('id');
    $table->string('owner');
    $table->decimal('balance', 10, 2)->default(0);
});

// 初始化账户
Capsule::table('accounts')->insert([
    ['owner' => 'Alice', 'balance' => 1000.00],
    ['owner' => 'Bob',   'balance' => 500.00],
]);

function printBalances() {
    $accounts = Capsule::table('accounts')->orderBy('id')->get();
    foreach ($accounts as $a) {
        echo "  {$a->owner}: ¥" . number_format($a->balance, 2) . "\n";
    }
}

echo "=== 初始余额 ===\n";
printBalances();

echo "\n=== 事务：成功提交 ===\n";
try {
    $conn = Capsule::connection();
    $conn->beginTransaction();

    // Alice 转账 200 给 Bob
    Capsule::table('accounts')->where('owner', 'Alice')->decrement('balance', 200);
    Capsule::table('accounts')->where('owner', 'Bob')->increment('balance', 200);

    // 检查余额是否足够（这里余额充足，事务应提交）
    $alice = Capsule::table('accounts')->where('owner', 'Alice')->first();
    if ($alice->balance < 0) {
        throw new Exception('余额不足，事务回滚');
    }

    $conn->commit();
    echo "✓ 转账成功，事务已提交\n";
} catch (Throwable $e) {
    echo "✗ 转账失败: " . $e->getMessage() . "\n";
}

echo "转账后余额:\n";
printBalances();

echo "\n=== 事务：失败回滚 ===\n";
try {
    $conn = Capsule::connection();
    $conn->beginTransaction();

    // Bob 尝试转账 10000 给 Alice（余额不足）
    Capsule::table('accounts')->where('owner', 'Bob')->decrement('balance', 10000);
    Capsule::table('accounts')->where('owner', 'Alice')->increment('balance', 10000);

    // 检查余额
    $bob = Capsule::table('accounts')->where('owner', 'Bob')->first();
    if ($bob->balance < 0) {
        throw new Exception('Bob 余额不足，事务回滚');
    }

    $conn->commit();
    echo "✓ 转账成功（不应出现）\n";
} catch (Throwable $e) {
    // 捕获异常后回滚事务
    $conn->rollBack();
    echo "✗ " . $e->getMessage() . "\n";
}

echo "\n回滚后余额（应恢复原状）:\n";
printBalances();

echo "\n=== 使用闭包方式执行事务 ===\n";
// Illuminate Database 也支持通过闭包执行事务
$result = Capsule::transaction(function () {
    Capsule::table('accounts')->where('owner', 'Alice')->increment('balance', 100);
    Capsule::table('accounts')->where('owner', 'Bob')->decrement('balance', 100);
    return 'transaction completed';
});

echo "闭包事务结果: {$result}\n";
echo "事务后余额:\n";
printBalances();

echo "\n✓ 事务示例执行成功\n";
