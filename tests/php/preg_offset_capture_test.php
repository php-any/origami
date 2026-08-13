<?php
$pattern = '#\{(!)?([\w\x80-\xFF]+)\}#';
$subject = '/';
preg_match_all($pattern, $subject, $matches, PREG_OFFSET_CAPTURE | PREG_SET_ORDER);
echo "count=" . count($matches) . "\n";
var_dump($matches);

$subject2 = '/users/{id}';
preg_match_all($pattern, $subject2, $matches2, PREG_OFFSET_CAPTURE | PREG_SET_ORDER);
echo "count2=" . count($matches2) . "\n";
var_dump($matches2);
