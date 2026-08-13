<?php

namespace tests\php;

require dirname(__DIR__, 2) . '/examples/laravel/vendor/autoload.php';

use Symfony\Component\String\UnicodeString;

$u = new UnicodeString(' [OK] Formatted success message');
Log::info('length=' . $u->length());
Log::info('width=' . $u->width(false));

// OutputWrapper
use Symfony\Component\Console\Helper\OutputWrapper;
$w = new OutputWrapper();
$wrapped = $w->wrap('Formatted success message', 70, "\n");
Log::info('wrapped: ' . var_export($wrapped, true));

Log::info('done');
