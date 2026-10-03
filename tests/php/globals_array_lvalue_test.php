<?php
$sharedGlobal = ['value' => 1];
if (!is_array($GLOBALS) || $GLOBALS['sharedGlobal'] !== $sharedGlobal) { throw new Exception('global snapshot kind'); }
$snapshot = $GLOBALS;
$GLOBALS['sharedGlobal']['value'] = 2;
if ($sharedGlobal !== ['value' => 2] || $snapshot['sharedGlobal'] !== ['value' => 1]) { throw new Exception('global dimension COW'); }
$reference =& $GLOBALS['sharedGlobal']['value'];
$reference = 3;
if ($sharedGlobal['value'] !== 3) { throw new Exception('global dimension reference'); }
$GLOBALS['fromDimension'] = 4;
function readDimensionGlobal() { global $fromDimension; return $fromDimension; }
if (readDimensionGlobal() !== 4) { throw new Exception('dimension to global binding'); }
unset($GLOBALS['fromDimension']);
if (isset($GLOBALS['fromDimension'])) { throw new Exception('global unset'); }
$_SESSION = ['value' => 5];
$sessionCopy = $_SESSION;
$_SESSION['value'] = 6;
if (!is_array($_SESSION) || $sessionCopy !== ['value' => 5]) { throw new Exception('session COW'); }
echo "global and session arrays OK\n";
