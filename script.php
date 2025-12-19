<?php

$a = 1;

beginning:
if($a == 2) {
    echo "ok";
    return;
}
$a = 2;
goto beginning;

