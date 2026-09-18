<?php

require __DIR__.'/../../vendor/autoload.php';

class CollectNotCallerClass_Pkg
{
    public function hasCommands(...$commandClassNames)
    {
        return collect($commandClassNames)->flatten()->toArray();
    }
}

$p = new CollectNotCallerClass_Pkg();
$got = $p->hasCommands(['A', 'B'], 'C');
if ($got !== ['A', 'B', 'C']) {
    fwrite(STDERR, 'collect() 在类方法中应返回 Collection 而非调用方类: '.json_encode($got)."\n");
    exit(1);
}
echo "collect_not_caller_class ok\n";
