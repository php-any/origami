<?php

function on_optional($value) {
    return new class ($value) {
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

$a = ['v' => on_optional(null)->getActionName()];
echo "ok\n";
