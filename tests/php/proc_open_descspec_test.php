<?php
namespace tests\php;

$descriptorspec = [
    1 => ['pipe', 'w'],
    2 => ['pipe', 'w'],
];
Log::info('keys=' . json_encode(array_keys($descriptorspec)));
foreach ($descriptorspec as $k => $v) {
    Log::info("k=$k type=" . get_debug_type($v) . ' v0=' . $v[0] . ' v1=' . $v[1]);
}

$pipes = null;
$process = @proc_open(['stty', '-a'], $descriptorspec, $pipes, null, null, ['suppress_errors' => true]);
Log::info('process=' . get_debug_type($process) . ' false=' . ($process === false ? '1' : '0'));
Log::info('pipes=' . get_debug_type($pipes));
if (is_array($pipes) || is_object($pipes)) {
    foreach ($pipes as $k => $v) {
        Log::info("pipe[$k]=" . get_debug_type($v));
    }
}
if ($process && isset($pipes[1])) {
    $info = stream_get_contents($pipes[1]);
    Log::info('sgc type=' . get_debug_type($info) . ' is_false=' . ($info === false ? '1' : '0'));
}

// Mimic exact return
function readFromProcess_probe($command): ?string {
    $descriptorspec = [
        1 => ['pipe', 'w'],
        2 => ['pipe', 'w'],
    ];
    $pipes = [];
    if (!$process = @proc_open($command, $descriptorspec, $pipes, null, null, ['suppress_errors' => true])) {
        return null;
    }
    $info = stream_get_contents($pipes[1]);
    Log::info('inside info type=' . get_debug_type($info));
    fclose($pipes[1]);
    fclose($pipes[2]);
    proc_close($process);
    return $info;
}

try {
    $r = readFromProcess_probe(['stty', '-a']);
    Log::info('probe result type=' . get_debug_type($r) . ' len=' . (is_string($r) ? strlen($r) : -1));
} catch (\Throwable $e) {
    Log::fatal('probe threw: ' . $e->getMessage());
}

Log::info('proc_open_descspec_test 测试通过');
