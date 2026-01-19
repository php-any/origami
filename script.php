<?php

class A {
    public function __construct(private string $name = "A", $a = "b")
    {
        var_dump($name);
        var_dump($a);
    }

    public function getName():string {
        return $this->name;
    }
}

class B extends A
{
    public function __construct()
    {
        parent::__construct("name", "d2");
    }
}

$data = new B();
var_dump($data);