<?php
// PHP: 键存在且值为 null 时，读取不 Warning，isset 为 false
$a = [0 => null];
var_dump($a[0]);
var_dump(isset($a[0]));
$x = $a[0] ?? 'default';
var_dump($x);
