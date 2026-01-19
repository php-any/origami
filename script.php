<?php

class AA {
    public function __construct($a = "aa")
    {
        var_dump($a);
    }
}

class A extends AA{
    public function __construct(private string $name = "A", $a = "b")
    {
        var_dump($name);
        var_dump($a);
        parent::__construct($a);
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