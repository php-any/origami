<?php
echo "p1\n";
var_dump(parse_url('http://127.0.0.1:8000'));
echo "p2\n";
var_dump(parse_url('http://127.0.0.1:8000/admin/login'));
echo "p3\n";
$bad = 'http://' . chr(0) . 'host';
var_dump(@parse_url($bad));
echo "done\n";
