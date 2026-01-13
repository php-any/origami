<?php
#[Deprecated]
class TestClass {
}
$reflection = new ReflectionClass(TestClass::class);
var_dump($reflection);
