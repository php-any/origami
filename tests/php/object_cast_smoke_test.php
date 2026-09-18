<?php
namespace tests\php;
$a = (object) ['x' => 1];
\Log::info('type='.gettype($a).' class='.(is_object($a)?get_class($a):'-'));
$b = (object) array('y' => 2);
\Log::info('b type='.gettype($b));
