<?php

if (getenv('ORIGAMI_WAIT_SIGNAL') !== '1') {
    Log::info('skip: shutdown_on_signal 需手动发送 SIGTERM/SIGINT，批量测试中跳过');
    return;
}

register_shutdown_function(function () {
    echo "SHUTDOWN_OK\n";
});

echo "READY\n";

// 阻塞直到收到退出信号（内置 handler 会同步执行 RunShutdownCallbacks）
Signal\wait([SIGTERM, SIGINT]);
