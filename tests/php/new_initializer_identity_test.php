<?php
class InitializerIdentityItem { public int $value = 0; }
function initializerDefault($item = new InitializerIdentityItem) { return $item; }
function initializerStatic() {
    static $item = new InitializerIdentityItem;
    return $item;
}
$first = initializerDefault();
$second = initializerDefault();
$first->value = 7;
if ($first === $second || $second->value !== 0) {
    throw new Exception('default new expression must create a fresh object per call');
}
$shared = initializerStatic();
$shared->value = 9;
if (initializerStatic() !== $shared || initializerStatic()->value !== 9) {
    throw new Exception('static new initializer must run once');
}
echo "new initializer identity OK\n";
