<?php
require __DIR__ . '/laravel/vendor/autoload.php';
$s = Termwind\Components\Span::fromStyles(null, 'test');
echo get_class($s) . "\n";
