<?php

namespace tests\php;

$a = [];
echo "type=".gettype($a)."\n";
echo "is_array=".var_export(is_array($a), true)."\n";
echo "bool=".var_export((bool)$a, true)."\n";
echo "count=".count($a)."\n";
echo "empty=".var_export(empty($a), true)."\n";

$b = [1];
echo "b_bool=".var_export((bool)$b, true)."\n";
