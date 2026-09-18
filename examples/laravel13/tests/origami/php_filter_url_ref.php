<?php
$cases = [
    'js/x',
    '/js/x',
    'http://a.com/x',
    'https://a.com',
    'ftp://a.com',
    'example.com',
    'http://',
    'http://a',
    'mailto:a@b.com',
    '//cdn.com/x',
    'http://127.0.0.1:8000/js/x',
];
foreach ($cases as $p) {
    echo $p, ' => ', var_export(filter_var($p, FILTER_VALIDATE_URL), true), PHP_EOL;
}
echo 'const=', FILTER_VALIDATE_URL, PHP_EOL;
