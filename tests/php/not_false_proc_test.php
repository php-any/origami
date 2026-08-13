<?php
namespace tests\php;

$v = false;
if (!$v) {
    Log::info('!false ok');
} else {
    Log::fatal('!false should be true');
}

$p = false;
if (!$p = false) {
    Log::info('!assign false ok');
} else {
    Log::fatal('!assign false should enter if');
}

// mimic: if (!$process = proc_open_fail)
$descriptorspec = [1 => ['pipe', 'w'], 2 => ['pipe', 'w']];
// force fail with empty command array?
$pipes = [];
$r = @proc_open([], $descriptorspec, $pipes);
Log::info('empty array proc=' . var_export($r, true) . ' type=' . get_debug_type($r));
if (!$process = @proc_open([], $descriptorspec, $pipes)) {
    Log::info('entered if on failed proc_open');
} else {
    Log::fatal('should have entered if, process type=' . get_debug_type($process));
}

Log::info('not_false_proc_test 测试通过');
