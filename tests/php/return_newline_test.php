<?php

class A {
    public function withNewline() {
        $x = 'ok';
        return
            $x;
    }
    public function withSameLine() {
        $x = 'ok';
        return $x;
    }
    public function chainNewline() {
        $arr = [1, 2, 3];
        return
            $arr;
    }
}

$a = new A();
echo "newline: " . var_export($a->withNewline(), true) . "\n";
echo "sameline: " . var_export($a->withSameLine(), true) . "\n";
echo "chain: " . var_export($a->chainNewline(), true) . "\n";
