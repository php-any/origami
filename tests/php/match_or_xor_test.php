<?php

$aNeg = true;
$bNeg = false;
$operator = 'or';
$negative = match ($operator) {
    'and' => $aNeg and $bNeg,
    'or' => $aNeg or $bNeg,
    'xor' => $aNeg xor $bNeg,
};
var_dump($negative);

$operator = 'xor';
$negative = match ($operator) {
    'and' => $aNeg and $bNeg,
    'or' => $aNeg or $bNeg,
    'xor' => $aNeg xor $bNeg,
};
var_dump($negative);
