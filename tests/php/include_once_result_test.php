<?php
$first = require_once __DIR__ . '/fixtures/include_once_result.php';
$second = require_once __DIR__ . '/fixtures/include_once_result.php';
if ($first !== ['value' => 1] || $second !== true) { throw new Exception('require_once return semantics'); }
$third = include __DIR__ . '/fixtures/include_once_result.php';
if ($third !== ['value' => 1]) { throw new Exception('ordinary include must execute again'); }
echo "include_once result OK\n";
