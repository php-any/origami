<?php
namespace NativeInjectRegression;

final class Dependency { public string $value = 'native'; }
final class Ordinary {
    public Dependency $dependency;
    public function __construct(#[\Container\Annotation\Inject('native.probe')] Dependency $dependency) { $this->dependency = $dependency; }
}
final class Promoted {
    public function __construct(#[\Container\Annotation\Inject('native.probe')] public Dependency $dependency) {}
}
class NamedParent {
    public Dependency $dependency;
    public function __construct(#[\Container\Annotation\Named('native.probe')] Dependency $dependency) { $this->dependency = $dependency; }
}
final class NamedChild extends NamedParent {}
final class Inferred {
    public Dependency $dependency;
    public function __construct(#[\Container\Annotation\Inject] Dependency $dependency) { $this->dependency = $dependency; }
}

$container = new \Container\Container();
$dependency = new \NativeInjectRegression\Dependency();
$inferred = new \NativeInjectRegression\Dependency();
$container->instance('native.probe', $dependency);
$container->instance(\NativeInjectRegression\Dependency::class, $inferred);
foreach ([\NativeInjectRegression\Ordinary::class, \NativeInjectRegression\Promoted::class, \NativeInjectRegression\NamedParent::class, \NativeInjectRegression\NamedChild::class] as $class) {
    if ($container->make($class)->dependency !== $dependency) { throw new \RuntimeException('Native constructor injection failed: ' . $class); }
}
if ($container->make(\NativeInjectRegression\Inferred::class)->dependency !== $inferred) { throw new \RuntimeException('Native Inject type inference failed'); }
echo "native constructor annotations OK\n";
