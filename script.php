<?php

class Foo {
    public function handleException($e) {
        echo "handleException called: " . $e . "\n";
    }

    public function bar() {
        $method = "handleException";
        $this->{$method}("test error");
    }
}

$obj = new Foo();
$method = "handleException";
$obj->{$method}("hello");
$obj->bar();
