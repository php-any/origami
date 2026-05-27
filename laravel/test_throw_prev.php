<?php

require __DIR__.'/vendor/autoload.php';

function inner() {
    try {
        new ReflectionClass('auth');
    } catch (ReflectionException $e) {
        throw new Illuminate\Contracts\Container\BindingResolutionException("Target class [auth] does not exist.", 0, $e);
    }
}

function outer() {
    try {
        inner();
    } catch (Throwable) {
        echo "caught\n";
        return [];
    }
}

var_export(outer());
