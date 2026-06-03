<?php
class TestLib {
    public $name;

    public function __construct($name) {
        $this->name = $name;
    }

    public function greet() {
        return "Hello, " . $this->name;
    }
}

function test_add($a, $b) {
    return $a + $b;
}
