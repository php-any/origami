<?php
preg_match_all('/<x>/', 'plain text', $matches, PREG_OFFSET_CAPTURE);
var_dump($matches);
foreach ($matches[0] as $i => $match) {
    echo "iter\n";
}
