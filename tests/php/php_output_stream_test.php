<?php
ob_start();
$out = fopen('php://output', 'w');
if (!is_resource($out)) { throw new Exception('output resource'); }
$n = fwrite($out, 'hello');
echo ':' . $n . ':';
fclose($out);
$s = ob_get_clean();
if ($s !== "hello:5:") { throw new Exception('output stream: ' . $s); }
echo "OK\n";
