<?php
namespace ComposerCaseProbe;
class Box {}
$object = new Box();
$matched = false;
switch (true) {
case $object instanceof Box:
    $matched = true;
    break;
default:
    throw new \Exception('Wrong instanceof case');
}
if (!$matched) {
    throw new \Exception('Class name parsed as label');
}
echo "instanceof case passed\n";
