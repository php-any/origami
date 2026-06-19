<?php
function f() {
  static $n = 0;
  echo $n++;
  echo "\n";
}
f();
f();
f();
