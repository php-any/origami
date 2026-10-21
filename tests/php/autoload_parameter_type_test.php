<?php
namespace ComposerTypeProbe;
function acceptContract(Contract $value) {
    return $value->value();
}
spl_autoload_register(function ($name) {
    if ($name === __NAMESPACE__ . '\\Contract') {
        eval('namespace ComposerTypeProbe; interface Contract { public function value(); }');
    }
});
class Implementation implements Contract {
    public function value() { return 42; }
}
if (acceptContract(new Implementation()) !== 42) {
    throw new \Exception('Unloaded parameter interface was not namespace-qualified');
}
echo "autoload parameter type passed\n";
