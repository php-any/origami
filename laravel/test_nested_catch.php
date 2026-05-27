<?php

require __DIR__.'/vendor/autoload.php';

function simulateBuild($concrete) {
    try {
        $reflector = new ReflectionClass($concrete);
    } catch (ReflectionException $e) {
        throw new Illuminate\Contracts\Container\BindingResolutionException("Target class [$concrete] does not exist.", 0, $e);
    }
    return $reflector;
}

function outer() {
    try {
        simulateBuild('auth');
    } catch (Throwable) {
        echo "outer caught\n";
        return [];
    }
}

var_export(outer());
