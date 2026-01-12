<?php

class Demo {
    public function __construct(
        private string|int|null $name = null,
        public string|int|bool $value = 0
    ) {
        $this->name = $name;
    }
    
    public function getName() {
        return $this->name;
    }
    
    public function getValue() {
        return $this->value;
    }
}

$demo1 = new Demo(11);
echo "Demo1 - name: ";
var_dump($demo1->getName());


$demo1->value = 222;
var_dump($demo1->getValue());
