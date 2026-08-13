<?php

$env = getenv();
if (!is_array($env)) {
    Log::fatal('getenv() 无参数应返回 array');
}

$merged = $_ENV + $env;
if (!is_array($merged)) {
    Log::fatal('$_ENV + getenv() 应返回 array');
}

Log::info('getenv_all_test OK');
