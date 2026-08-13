<?php
namespace tests\php;

/**
 * proc_open 数组命令 + 连续调用（ProgressBar / Terminal::readFromProcess）。
 */

$descriptorspec = [1 => ['pipe', 'w'], 2 => ['pipe', 'w']];
$pipes = [];
$p = proc_open(['/bin/echo', 'hello-stty'], $descriptorspec, $pipes);
$out = stream_get_contents($pipes[1]);
fclose($pipes[1]);
fclose($pipes[2]);
proc_close($p);
Log::info('array echo out=' . var_export($out, true));
if (trim((string) $out) !== 'hello-stty') {
    Log::fatal('array command should run without shell, got ' . var_export($out, true));
}

for ($i = 0; $i < 30; $i++) {
    $pipes = [];
    $p = @proc_open(['stty', '-a'], $descriptorspec, $pipes, null, null, ['suppress_errors' => true]);
    if ($p === false) {
        Log::fatal("proc_open false at $i");
    }
    $info = stream_get_contents($pipes[1]);
    if ($info === false) {
        Log::fatal("stream_get_contents false at $i");
    }
    if (!is_string($info)) {
        Log::fatal("not string at $i: " . get_debug_type($info));
    }
    fclose($pipes[1]);
    fclose($pipes[2]);
    proc_close($p);
}
Log::info('proc_open_array_cmd_test 测试通过');
