<?php
namespace test

class LoadClass {
    public $map = [];

    public function load($class)
    {
        var_dump($this->map);
    }
}


$data = new LoadClass();

$data->map = [
    "test" => "test.php"
];

spl_autoload_register(array($data, 'load'));

new Test();

