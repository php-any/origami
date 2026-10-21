<?php
$max = PHP_INT_MAX . '.0.0.0';
if (version_compare($max, '2.0.0.0-dev') !== 1 || version_compare('0.0.0.0-dev', $max) !== -1) {
    throw new Exception('Large version component overflow');
}
$versions = ['1.0.0.0-dev', $max, '0.0.0.0-dev', '2.0.0.0-dev'];
usort($versions, 'version_compare');
if ($versions !== ['0.0.0.0-dev', '1.0.0.0-dev', '2.0.0.0-dev', $max]) { throw new Exception('Version interval sorting failed'); }
echo "large version comparison passed\n";
