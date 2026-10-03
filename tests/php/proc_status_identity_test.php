<?php
$command = PHP_OS_FAMILY === 'Windows'
    ? ['ping.exe', '-n', '3', '127.0.0.1']
    : ['sh', '-c', 'sleep 2'];
$process = proc_open($command, [1 => ['pipe', 'w'], 2 => ['pipe', 'w']], $pipes);
if (!is_resource($process)) {
    throw new Exception('proc_open failed');
}
try {
    $status = proc_get_status($process);
    if (!is_array($status) || is_object($status) || gettype($status) !== 'array') {
        throw new Exception('process status array identity');
    }
    if ($status['running'] !== true || $status['exitcode'] !== -1 || $status['pid'] <= 0) {
        throw new Exception('live process was reported as exited');
    }
    $copy = $status;
    $copy['running'] = false;
    if ($status['running'] !== true || proc_get_status($process)['running'] !== true) {
        throw new Exception('status query changed process state');
    }
} finally {
    proc_close($process);
    foreach ($pipes as $pipe) {
        if (is_resource($pipe)) {
            throw new Exception('proc_close left an open pipe');
        }
    }
}
if (is_resource($process)) {
    throw new Exception('proc_close left an open process resource');
}
echo "Process status identity OK\n";
