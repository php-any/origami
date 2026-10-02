<?php
function check_resource($actual, $expected, $label) {
    if ($actual !== $expected) { throw new Exception($label . ': ' . json_encode($actual)); }
}
$path = __DIR__ . '/../../examples/laravel13/storage/origami-debug/request-resource-' . getmypid() . '.txt';
$copy = $path . '.copy';
try {
    check_resource(file_put_contents($path, "first\n"), 6, 'write');
    check_resource(file_put_contents($path, "second\n", FILE_APPEND), 7, 'append');
    check_resource(file_get_contents($path), "first\nsecond\n", 'read whole file');
    check_resource(file($path, FILE_IGNORE_NEW_LINES), ['first', 'second'], 'read lines');
    check_resource(md5_file($path), md5("first\nsecond\n"), 'md5 file');
    check_resource(hash_file('sha256', $path), hash('sha256', "first\nsecond\n"), 'hash file');
    check_resource(copy($path, $copy), true, 'copy');
    check_resource(file_get_contents($copy), "first\nsecond\n", 'copy content');
    $stream = fopen($path, 'r');
    check_resource(fread($stream, 6), "first\n", 'read stream');
    check_resource(stream_get_contents($stream), "second\n", 'read remaining');
    check_resource(fclose($stream), true, 'close');
    check_resource(shell_exec('echo ready'), "ready\n", 'shell text output');
    check_resource(shell_exec('exit 0'), null, 'shell no output');
    $command = PHP_OS_FAMILY === 'Windows' ? [getenv('COMSPEC'), '/C', 'echo', 'proc-output'] : ['/bin/sh', '-c', 'echo proc-output'];
    $process = proc_open($command, [1 => ['pipe', 'w'], 2 => ['pipe', 'w']], $pipes);
    check_resource(is_array($pipes), true, 'proc pipes are array');
    check_resource(array_keys($pipes), [1, 2], 'proc pipe keys');
    check_resource(str_replace("\r\n", "\n", stream_get_contents($pipes[1])), "proc-output\n", 'proc output');
    fclose($pipes[1]);
    fclose($pipes[2]);
    check_resource(proc_close($process), 0, 'proc close');
    $caught = false;
    try { usleep(-1); } catch (ValueError $e) { $caught = true; }
    check_resource($caught, true, 'usleep negative value');
    check_resource(usleep(0), null, 'usleep zero');
    echo "request resource basic PASS\n";
} finally {
    unlink($path);
    unlink($copy);
}
