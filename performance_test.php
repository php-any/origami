<?php

echo "=== 性能测试：一百万次赋值 ===\n";

// 记录开始时间
$startTime = microtime(true);

// 执行一百万次赋值操作
for ($i = 1; $i <= 1000000; $i++) {
    $value = $i;
    $result = $value * 2;
    $sum = $value + $result;
}

// 记录结束时间
$endTime = microtime(true);

// 计算执行时间
$executionTime = $endTime - $startTime;

echo "执行完成！\n";
echo "总执行时间: " . number_format($executionTime, 4) . " 秒\n";
echo "平均每次操作时间: " . number_format($executionTime / 1000000 * 1000000, 6) . " 微秒\n";
echo "每秒操作次数: " . number_format(1000000 / $executionTime, 0) . " 次/秒\n";


echo "时间信息：", $startTime, " -> ", $endTime, "\n";
echo "循环次数：", $i, "\n";
echo "循环内 $value ：", $value, "\n";
echo "循环内 $result ：", $result, "\n";
echo "循环内 $sum ：", $sum, "\n";

echo "\n=== 性能测试完成 ===\n"; 