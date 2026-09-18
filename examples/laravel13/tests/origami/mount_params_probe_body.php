<?php
/** @var mixed $app */

echo "step1 anon class get_class/class_exists\n";
$obj = new class extends \Livewire\Component {
    public $forMount;
};
$cn = get_class($obj);
echo "class_name=$cn\n";
echo "class_exists=".(class_exists($cn)?'y':'n')."\n";
echo "is_subclass=".(is_subclass_of($cn, \Livewire\Component::class)?'y':'n')."\n";
try {
    $o2 = new $cn;
    echo "new_by_name_ok class=".get_class($o2)."\n";
} catch (Throwable $e) {
    echo "new_by_name EX: ".$e->getMessage()."\n";
}

echo "step2 livewire component register\n";
try {
    app('livewire')->component('__mountParamsContainer', $obj);
    echo "register_ok\n";
    $exists = app('livewire')->exists('__mountParamsContainer');
    echo "exists=".($exists?'y':'n')."\n";
    $c = app('livewire')->new('__mountParamsContainer');
    echo "new_ok class=".get_class($c)."\n";
} catch (Throwable $e) {
    echo "livewire EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
}

echo "step3 SupportLazyLoading registerContainerComponent\n";
try {
    $hook = new \Livewire\Features\SupportLazyLoading\SupportLazyLoading();
    $ref = new ReflectionClass($hook);
    $m = $ref->getMethod('registerContainerComponent');
    $m->setAccessible(true);
    $m->invoke($hook);
    echo "registerContainer_ok\n";
    $c = app('livewire')->new('__mountParamsContainer');
    echo "new_after_hook_ok class=".get_class($c)."\n";
} catch (Throwable $e) {
    echo "hook EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
}

echo "DONE\n";
