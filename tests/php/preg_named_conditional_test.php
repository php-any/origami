<?php
$pattern = '#^(\$|(?P<percent>%))(?P<var>\w++)(?(percent)%)(?P<path>.*)#';
foreach (['package' => 'package', '$HOME/src' => 'HOME:/src', '%HOME%/src' => 'HOME:/src'] as $input => $expected) {
    $actual = preg_replace_callback($pattern, function ($m) { return $m['var'] . ':' . $m['path']; }, $input);
    if ($actual !== $expected) { throw new Exception('Named conditional replacement failed: ' . $input); }
}
echo "named conditional regex passed\n";
