<?php

$a = ['x' => null];
var_dump(array_key_exists('x', $a));
var_dump(isset($a['x']));
Log::info('key exists test');
