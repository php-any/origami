<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Testing\AssertableJsonString;

/**
 * 不调用依赖 PHPUnit 的 assert* 方法；仅验证 JSON 解码与 data_get 路径。
 */
$json = new AssertableJsonString('{"user":{"id":1,"name":"Origami"}}');
if ($json->json('user.name') !== 'Origami') {
    echo "FAIL: json path\n";
    var_export($json->json('user.name'));
    echo "\n";
    exit(1);
}

if (count($json) !== 1) {
    echo "FAIL: count\n";
    exit(1);
}

if (!isset($json['user'])) {
    echo "FAIL: array access\n";
    exit(1);
}

echo "PASS\n";
