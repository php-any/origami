<?php
$contents = file_get_contents(__FILE__);
$file = new SplFileObject(__FILE__, 'r');
if ($file->ftell() !== 0 || $file->fread(4) !== substr($contents, 0, 4) || $file->ftell() !== 4) { throw new Exception('byte cursor'); }
if ($file->fseek(-2, SEEK_CUR) !== 0 || $file->fread(3) !== substr($contents, 2, 3)) { throw new Exception('relative byte seek'); }
if ($file->fseek(-3, SEEK_END) !== 0 || $file->fread(9) !== substr($contents, -3) || !$file->eof()) { throw new Exception('byte EOF'); }
if ($file->fseek(0) !== 0 || $file->eof() || $file->fread(5) !== substr($contents, 0, 5)) { throw new Exception('seek clears EOF'); }
echo "SPL byte cursor OK\n";
