<?php
set_error_handler(function ($severity, $message) { throw new Exception($message); });
trait InheritedRequirement { abstract public function value(): int; }
class InheritedImplementation {
    private $stored;
    public function value(): int { $this->stored = 42; return $this->stored; }
}
class InheritedUser extends InheritedImplementation { use InheritedRequirement; }
if ((new InheritedUser())->value() !== 42) { throw new Exception('Inherited implementation lost lexical scope'); }
restore_error_handler();
echo "abstract trait inherited scope passed\n";
