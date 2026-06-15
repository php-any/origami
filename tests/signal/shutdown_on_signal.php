<?php

register_shutdown_function(function () {
    echo "SHUTDOWN_OK\n";
});

echo "READY\n";

// 阻塞直到收到退出信号（内置 handler 会同步执行 RunShutdownCallbacks）
Signal\wait([SIGTERM, SIGINT]);
