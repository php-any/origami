<?php

$values = ['b' => 20, 'a' => 10, 'c' => 30];
asort($values);
if (array_keys($values) !== ['a', 'b', 'c']) {
    Log::fatal('asort 未按值升序并保留键');
}

arsort($values);
if (array_keys($values) !== ['c', 'b', 'a']) {
    Log::fatal('arsort 未按值降序并保留键');
}

Log::info('asort_arsort_test OK');
