<?php
declare(strict_types=1);
require __DIR__.'/../../vendor/autoload.php';

$svg = new BladeUI\Icons\Svg('callback-regression', '<svg></svg>', [
    'enabled' => true,
    'disabled' => false,
    'tabindex' => 12,
]);
if ($svg->toHtml() !== '<svg enabled="1" disabled="" tabindex="12"></svg>') {
    throw new RuntimeException('Blade icon callback scalar conversion failed');
}
return true;
