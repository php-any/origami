<?php

class Demo {
    public function __construct(
        private string|int|null $name = null,
        public string|int|bool $value = 0,
        string|int $mixed = 'default'
    ) {
    }
}

$demo1 = new Demo();
dump($demo1);

$demo2 = new Demo('test', 123, 456);
dump($demo2);

$demo3 = new Demo(null, true, 'string');
dump($demo3);

