<?php

namespace tests\php;

$prev = set_error_handler(function () {});
restore_error_handler();
Log::info('restore_error_handler 测试通过');
