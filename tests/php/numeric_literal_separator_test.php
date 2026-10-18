<?php
$values = [1_234_567, 0xCA_FE, 0b1010_0101, 0o7_55, 07_55, 1_2.5_0, 1_2e1_0];
if ($values !== [1234567, 51966, 165, 493, 493, 12.5, 120000000000.0]) {
    throw new Exception('numeric separators must preserve literal values and types');
}
echo "numeric literal separators OK\n";
