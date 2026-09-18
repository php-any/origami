<?php
// top-level short echo mix (no class)
ob_start();
echo 'A';
?><?= 'CHILD_HTML' ?><?php
echo 'B';
$out = ob_get_clean();
if ($out !== 'ACHILD_HTMLB') {
    fwrite(STDERR, 'FAIL top: '.var_export($out, true)."\n");
    exit(1);
}
echo "top_ok\n";
