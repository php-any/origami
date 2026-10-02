<?php
function weakString(string $value) { return $value; }
function weakReturn(): string { return 42; }
function weakCallStrict() { return strictString(42); }
function weakWrite($object) { return $object->text = 42; }
function weakMakeClosure() { return function (): string { return 42; }; }
class WeakTypedObject {
    public string $text = 'before';
    public function accepts(string $value) { return $value; }
    public function converted(): string { return 42; }
}
