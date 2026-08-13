<?php

function static_increment_f() {
  static $n = 0;
  echo $n++;
  echo "\n";
}
static_increment_f();
static_increment_f();
static_increment_f();
