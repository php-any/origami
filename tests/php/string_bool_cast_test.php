<?php

if ((bool) '') {
    echo "FAIL empty string\n";
    exit(1);
}

if ((bool) '0') {
    echo "FAIL zero string\n";
    exit(1);
}

if (!(bool) '00') {
    echo "FAIL non-zero string\n";
    exit(1);
}

if (!(bool) 'value') {
    echo "FAIL regular string\n";
    exit(1);
}

echo "PASS\n";
