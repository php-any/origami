<?php
$values = ["OrigamiFixture\\", "a\\'b", "nul\0byte", ['' => 'empty', "a\\" => "b\\"]];
foreach ($values as $value) {
    $code = var_export($value, true);
    if (eval('return ' . $code . ';') !== $value) { throw new Exception('Export does not round trip'); }
}
ob_start(); var_export('buffered'); $out = ob_get_clean();
if ($out !== "'buffered'") { throw new Exception('Export bypassed output buffer'); }
echo "var_export escaping passed\n";
