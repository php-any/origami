<?php

// 测试基本的 foreach 循环
$fruits = ["apple", "banana", "orange"];

foreach ($fruits as $fruit) {
    echo "水果: " . $fruit . "\n";
}

// 测试带键的 foreach 循环
$colors = ["red" => "红色", "green" => "绿色", "blue" => "蓝色"];

foreach ($colors as $key => $value) {
    echo "键: " . $key . ", 值: " . $value . "\n";
}

// 测试数字数组的 foreach 循环
$numbers = [1, 2, 3, 4, 5];

foreach ($numbers as $index => $number) {
    echo "索引: " . $index . ", 数字: " . $number . "\n";
}

// 测试嵌套数组
$matrix = [
    [1, 2, 3],
    [4, 5, 6],
    [7, 8, 9]
];

foreach ($matrix as $row_index => $row) {
    echo "第 " . $row_index . " 行: ";
    foreach ($row as $col_index => $value) {
        echo $value . " ";
    }
    echo "\n";
} 