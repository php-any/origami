<?php
trait ConstructorPrivateHelper { private function prepare(): void { $this->value = 3; } }
trait ConstructorBody { public function __construct() { $this->prepare(); } }
class ConstructorParent { use ConstructorPrivateHelper, ConstructorBody; public int $value = 0; }
class ConstructorChild extends ConstructorParent {}
if ((new ConstructorChild())->value !== 3) { throw new Exception('inherited constructor lost its lexical scope'); }
echo "trait constructor visibility OK\n";
