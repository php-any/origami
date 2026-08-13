<?php

function ao_opt($x) {
    return new class ($x) {
        private $value;
        public function __construct($value) { $this->value = $value; }
        public function __call($method, $parameters) {
            if (is_object($this->value)) {
                return $this->value->{$method}(...$parameters);
            }
            return null;
        }
    };
}

$a = ['v' => ao_opt(null)->m()];
var_dump($a);
Log::info('anon explicit null ok');
