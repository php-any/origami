<?php
$path = tempnam(sys_get_temp_dir(), 'strip');
try {
    file_put_contents($path, '<?php /* comment */ class Example { /** doc */ public $text = "/* preserved */"; }');
    $actual = php_strip_whitespace($path);
    if ($actual !== '<?php  class Example { public $text = "/* preserved */"; }') {
        throw new Exception('Unexpected stripped source: ' . $actual);
    }
} finally { unlink($path); }
echo "php strip whitespace passed\n";
