<?php
namespace tests\php;
Log::info('LEFT=' . var_export(STR_PAD_LEFT, true) . ' type=' . get_debug_type(STR_PAD_LEFT));
Log::info('RIGHT=' . var_export(STR_PAD_RIGHT, true));
Log::info('BOTH=' . var_export(STR_PAD_BOTH, true));
$s = str_pad('a', 5, ' ', STR_PAD_LEFT);
Log::info('pad=' . var_export($s, true));
Log::info('done');
