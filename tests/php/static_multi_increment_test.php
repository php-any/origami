<?php
function f() {
  static $a = 0, $b = 0;
  echo $a++ . "," . $b++ . "\n";
}
f();
f();
f();
