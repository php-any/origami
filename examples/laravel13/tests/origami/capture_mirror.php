<?php
// Mirrors Filament SupportServiceProvider @capture registration pattern
namespace CaptureTest;

function register($store) {
    $store['capture'] = function (string $expression): string {
        [$name, $arguments] = str_contains($expression, ',') ?
            array_map('trim', explode(',', $expression, 2)) :
            [$expression, ''];

        return "
                <?php {$name} = (function (\$args) {
                    return function ({$arguments}) use (\$args) {
                        extract(\$args, EXTR_SKIP);
                        ob_start(); ?>
            ";
    };
}
