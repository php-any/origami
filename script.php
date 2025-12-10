<?php
namespace test

class Test {
    public static $a = "a";

    public function  b() {
        return self::$a;
    }
}

echo Test::$a;