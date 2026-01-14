<?php

class A {
    public function __construct(private string $name = "A", $a = "b")
    {
        echo $a;
    }

    public function getName():string {
        return $this->name;
    }
}

class B extends A
{

}

$data = new B(a: "c");
echo $data->getName()