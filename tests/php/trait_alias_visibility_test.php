<?php
trait VisibilityAliasTrait {
    public function value(): string { return 'value'; }
}
class VisibilityAliasTarget {
    use VisibilityAliasTrait { value as private hidden; }
    public function callHidden(): string { return $this->hidden(); }
}
$object = new VisibilityAliasTarget;
if ($object->value() !== 'value' || $object->callHidden() !== 'value') {
    throw new Exception('trait alias must preserve original and allow internal calls');
}
$caught = false;
try { $object->hidden(); } catch (Error $e) { $caught = true; }
if (!$caught || !(new ReflectionMethod($object, 'hidden'))->isPrivate()) {
    throw new Exception('trait alias visibility must be private');
}
echo "trait alias visibility OK\n";
