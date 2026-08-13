<?php
namespace tests\php;

class SpreadTyped_Event {
    public $name = 'x';
}

$listener = function (SpreadTyped_Event $event) {
    return $event->name;
};

$payload = [new SpreadTyped_Event()];
$r = $listener(...array_values($payload));
if ($r !== 'x') {
    Log::fatal('spread array_values 失败: '.var_export($r, true));
}

// also without array_values
$r2 = $listener(...$payload);
if ($r2 !== 'x') {
    Log::fatal('spread payload 失败: '.var_export($r2, true));
}

Log::info('spread_typed_closure 测试通过');
