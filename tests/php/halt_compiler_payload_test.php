<?php
$source = file_get_contents(__FILE__);
$payload = substr($source, __COMPILER_HALT_OFFSET__);
if (trim($payload) !== 'PAYLOAD: <?php invalid syntax {{{') {
    throw new Exception('compiler halt offset must point to uninterpreted payload');
}
echo "halt compiler payload OK\n";
__halt_compiler();
PAYLOAD: <?php invalid syntax {{{
