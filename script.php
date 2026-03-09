<?php

class User {
    public string $name;
}

$user = new User();
$user->name = "张三";

$data = [
    "a" => 1,
    "b" => 2,
    "c" => "ee",
    "user" => $user,
];

echo serialize($data);
