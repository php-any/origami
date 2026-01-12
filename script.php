<?php

class Demo {
    public function __construct(
        private string|int|null $name = null,
        public string|int|bool $value = 0
    ) {

    }
    
    public function getName() {
        return $this->name;
    }
    
    public function getValue() {
        return $this->value;
    }
}

$demo1 = new Demo(11, 2);
echo "Demo1 - name: ";
var_dump($demo1);
