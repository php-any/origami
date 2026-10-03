<?php
function fetch_identity_check($ok, $message) {
    if (!$ok) { throw new Exception($message); }
}
function fetch_identity_array(array $value): array { return $value; }
$pdo = new PDO('sqlite::memory:');
$query = "SELECT 'first' AS label, NULL AS missing, 'last' AS tail";
$assoc = $pdo->query($query)->fetch(PDO::FETCH_ASSOC);
fetch_identity_check(is_array($assoc) && !is_object($assoc) && gettype($assoc) === 'array', 'assoc identity');
fetch_identity_check(fetch_identity_array($assoc) === ['label' => 'first', 'missing' => null, 'tail' => 'last'], 'assoc order/values');
$copy = $assoc;
$copy['label'] = 'changed';
fetch_identity_check($assoc['label'] === 'first', 'assoc copy');
$both = $pdo->query($query)->fetch(PDO::FETCH_BOTH);
fetch_identity_check($both === ['label' => 'first', 0 => 'first', 'missing' => null, 1 => null, 'tail' => 'last', 2 => 'last'], 'both identity/keys');
$numeric = $pdo->query($query)->fetch(PDO::FETCH_NUM);
fetch_identity_check($numeric === ['first', null, 'last'], 'numeric result');
$object = $pdo->query($query)->fetch(PDO::FETCH_OBJ);
fetch_identity_check($object instanceof stdClass && is_object($object) && !is_array($object), 'object identity');
fetch_identity_check($object->label === 'first' && property_exists($object, 'missing') && $object->missing === null, 'object properties');
$all = $pdo->query($query)->fetchAll(PDO::FETCH_ASSOC);
fetch_identity_check($all === [$assoc], 'fetchAll array identity');
$objects = $pdo->query($query)->fetchAll(PDO::FETCH_OBJ);
fetch_identity_check(is_array($objects) && $objects[0] instanceof stdClass, 'fetchAll object identity');
echo "PDO fetch identity OK\n";
