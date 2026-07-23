<?php
namespace tests\php;

/**
 * 同一 $pipes 变量连续 proc_open（Symfony Terminal::readFromProcess 未每次重置 $pipes）。
 */

function read_twice(): ?string {
    $descriptorspec = [1 => ['pipe', 'w'], 2 => ['pipe', 'w']];
    // 故意不初始化 $pipes，并复用同一局部变量两次
    for ($i = 0; $i < 5; $i++) {
        if (!$process = @proc_open(['/bin/echo', "n$i"], $descriptorspec, $pipes, null, null, ['suppress_errors' => true])) {
            Log::fatal("proc_open failed i=$i");
        }
        $info = stream_get_contents($pipes[1]);
        if ($info === false) {
            Log::fatal("stream_get_contents false at i=$i");
        }
        fclose($pipes[1]);
        fclose($pipes[2]);
        proc_close($process);
        if (trim($info) !== "n$i") {
            Log::fatal("expected n$i got " . var_export($info, true));
        }
    }
    return $info;
}

read_twice();
Log::info('proc_open_reuse_pipes_test 测试通过');
