<?php
$source = fopen('php://temp', 'w+'); $target = fopen('php://temp', 'w+');
fwrite($source, 'abcdef'); rewind($source);
if (stream_copy_to_stream($source, $target) !== 6) { throw new Exception('Default stream copy failed'); }
rewind($target);
if (stream_get_contents($target) !== 'abcdef') { throw new Exception('Copied bytes differ'); }
if (stream_copy_to_stream($source, $target, 2, 1) !== 2) { throw new Exception('Offset copy failed'); }
if (stream_copy_to_stream($source, $target, 0) !== 0) { throw new Exception('Zero-length copy failed'); }
fclose($source); fclose($target);
echo "stream copy passed\n";
