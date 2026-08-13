<?php

namespace tests\php;

use Exception;

set_exception_handler(function ($e) {
    // 简单输出异常信息，验证是否进入回调
    Log::info('set_exception_handler 捕获到未处理异常: ' . $e->getMessage());
});

function throw_unhandled_exception() {
    throw new Exception('这是一个测试异常');
}

// 触发一个未被 try/catch 捕获的异常
throw_unhandled_exception();

