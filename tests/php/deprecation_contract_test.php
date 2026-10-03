<?php
$warnings = [];
set_error_handler(function ($level, $message) use (&$warnings) { $warnings[] = [$level, $message]; return true; });
class LegacyWarning implements Serializable {
    public function serialize(): string { return ''; }
    public function unserialize($data): void {}
}
class LegacyWarningChild extends LegacyWarning {}
class ModernWarning extends LegacyWarning {
    public function __serialize(): array { return []; }
    public function __unserialize(array $data): void {}
}
function deprecationCheck($ok, $message) { if (!$ok) throw new Exception($message); }
deprecationCheck(count($warnings) === 2 && $warnings[0][0] === E_DEPRECATED && str_contains($warnings[1][1], 'LegacyWarningChild implements'), 'Serializable declaration warning');
$warnings = [];
class DynamicWarning { public $declared = null; }
#[AllowDynamicProperties] class AllowedWarning {}
class AllowedWarningChild extends AllowedWarning {}
class StdWarningChild extends stdClass {}
$o = new DynamicWarning();
$o->declared = 1;
$o->created = 1;
$o->created = 2;
unset($o->created);
$o->created = 3;
$ref =& $o->reference;
$o->reference = 4;
foreach ([new stdClass(), new StdWarningChild(), new AllowedWarning(), new AllowedWarningChild()] as $allowed) $allowed->created = 1;
$copy = unserialize('O:14:"DynamicWarning":1:{s:7:"decoded";i:1;}');
deprecationCheck(count($warnings) === 4, 'property creation warning count');
deprecationCheck($warnings[0] === [E_DEPRECATED, 'Creation of dynamic property DynamicWarning::$created is deprecated'], 'property warning identity');
deprecationCheck($copy->decoded === 1 && $ref === 4, 'decoded and referenced properties');
set_error_handler(function ($level, $message) { throw new Exception('handler stopped'); });
try { $o->abort = 9; throw new Exception('missing warning'); }
catch (Exception $e) { deprecationCheck($e->getMessage() === 'handler stopped' && $o->abort === 9, 'write precedes warning and handler exception propagates'); }
restore_error_handler();
restore_error_handler();
enum NoDynamicWarning { case A; }
$case = NoDynamicWarning::A;
try { $case->extra = 1; throw new Exception('enum dynamic property accepted'); }
catch (Error $e) { deprecationCheck(!property_exists($case, 'extra'), 'enum write rejected before creation'); }
readonly class ReadonlyWarning {}
try { unserialize('O:15:"ReadonlyWarning":1:{s:5:"extra";i:1;}'); throw new Exception('readonly decoded dynamic property accepted'); }
catch (Error $e) { deprecationCheck(str_contains($e->getMessage(), 'Cannot create dynamic property'), 'readonly decode rejects creation'); }
echo "Deprecation contract OK\n";
