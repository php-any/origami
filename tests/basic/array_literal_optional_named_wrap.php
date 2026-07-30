<?php

class OnwWrap {
    public function __call($m, $a) {
        return null; // explicit
    }
}

function onw_optional($value) {
    return new OnwWrap();
}

$a = ['v' => onw_optional(null)->getActionName()];
var_dump($a);
Log::info('named wrap ok');
