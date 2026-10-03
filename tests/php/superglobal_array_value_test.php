<?php
function acceptsGlobalArray(array $value): bool { return is_array($value); }
foreach ([$_SERVER, $_GET, $_POST, $_COOKIE, $_FILES, $_ENV, $_REQUEST] as $global) {
    if (!acceptsGlobalArray($global)) { throw new Exception('superglobal array kind'); }
}
$_GET['key'] = 'first';
$copy = $_GET;
$_GET['key'] = 'second';
if ($copy['key'] !== 'first' || $_GET['key'] !== 'second') { throw new Exception('superglobal COW'); }
$reference =& $_GET;
$reference['key'] = 'third';
if ($_GET['key'] !== 'third') { throw new Exception('superglobal reference'); }
$_POST = 7;
if ($_POST !== 7) { throw new Exception('superglobal scalar assignment'); }
$_POST = null;
if ($_POST !== null) { throw new Exception('superglobal null assignment'); }
echo "superglobal array values OK\n";
