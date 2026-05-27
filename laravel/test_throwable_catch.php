<?php

require __DIR__.'/vendor/autoload.php';

function testThrowableCatch() {
    try {
        throw new Illuminate\Contracts\Container\BindingResolutionException('Target class [auth] does not exist.');
    } catch (Throwable) {
        echo "caught\n";
        return [];
    }
}

var_export(testThrowableCatch());
