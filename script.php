<?php

function gen_one_to_three() {
    yield 0;
    for ($i = 1; $i <= 4; $i++) {
        echo "gen_one_to_three start","\n";
        yield $i;
        echo "gen_one_to_three end","\n";
    }
    yield 5;
    yield 6;
}

$generator = gen_one_to_three();
foreach ($generator as $value) {
    echo "{$value}\n";
}