<?php
namespace tests\php;

/**
 * Symfony Terminal::readFromProcess 路径：proc_open(array) + stream_get_contents。
 */

if (!function_exists('proc_open')) {
    Log::fatal('proc_open missing');
}

$descriptorspec = [
    1 => ['pipe', 'w'],
    2 => ['pipe', 'w'],
];

$pipes = [];
$process = proc_open(['stty', '-a'], $descriptorspec, $pipes, null, null, ['suppress_errors' => true]);
if ($process === false) {
    Log::info('array cmd proc_open returned false, try string');
    $process = proc_open('stty -a', $descriptorspec, $pipes, null, null, ['suppress_errors' => true]);
}
if ($process === false) {
    Log::fatal('proc_open failed for both forms');
}

Log::info('pipes type=' . get_debug_type($pipes));
Log::info('isset pipes1=' . (isset($pipes[1]) ? '1' : '0'));
$info = stream_get_contents($pipes[1]);
Log::info('info type=' . get_debug_type($info) . ' len=' . (is_string($info) ? strlen($info) : -1));
fclose($pipes[1]);
fclose($pipes[2]);
proc_close($process);

if (!is_string($info)) {
    Log::fatal('stream_get_contents should return string, got ' . get_debug_type($info));
}

Log::info('proc_open_stty_test 测试通过');
