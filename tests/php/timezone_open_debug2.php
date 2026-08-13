<?php

namespace tests\php;

$tz = timezone_open('Europe/Paris');

Log::info('debug timezone_name_get raw: ' . timezone_name_get($tz));

