<?php

function onm_optional($value) {
    return new class ($value) {
        private $value;
        public function __construct($value) { $this->value = $value; }
        public function __call($method, $parameters) {
            if (is_object($this->value)) {
                return $this->value->{$method}(...$parameters);
            }
        }
    };
}

$a = ['v' => onm_optional(null)->getActionName()];
var_dump($a);
Log::info('optional null mimic ok');
