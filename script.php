<?php

class Option {
    private self $value;
    
    public function __construct(self $option) {
        $this->value = $option;
    }
    
    public function equals(self $option): bool {
        return $this->value === $option;
    }
}

$opt1 = new Option(new Option(null));
$opt2 = new Option(new Option(null));
var_dump($opt1->equals($opt2));
