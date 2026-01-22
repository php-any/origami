<?php

class A
{
    public $data = ["123"];
}

class B
{
    public function __construct(&$data)
    {
        $data = "456";
    }
}

$line = new A();;
$data = new B($line[0]);
var_dump($line);