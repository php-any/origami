<?php

function static_multi_increment_f() {
  static $a = 0, $b = 0;
  echo $a++ . "," . $b++ . "\n";
}
static_multi_increment_f();
static_multi_increment_f();
static_multi_increment_f();
