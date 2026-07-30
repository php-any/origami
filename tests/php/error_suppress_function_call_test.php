<?php

$length = @strlen('abc');
if ($length !== 3) {
    Log::fatal('@function() 未按错误抑制表达式执行');
}

Log::info('error_suppress_function_call_test OK');
