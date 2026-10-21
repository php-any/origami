<?php
class ConstantParentProbe {
    public const SHARED = 1;
    protected const PROTECTED_VALUE = 2;
    private const PRIVATE_VALUE = 3;
}
class ConstantChildProbe extends ConstantParentProbe {
    public const CHILD = 4;
}
$reflection = new \ReflectionClass(ConstantChildProbe::class);
$values = [];
foreach ($reflection->getReflectionConstants() as $constant) {
    if ($constant->isPublic()) {
        $values[$constant->getName()] = $constant->getValue();
    }
}
if ($values !== ['CHILD' => 4, 'SHARED' => 1]) {
    throw new \Exception('Inherited constant reflection failed');
}
$constant = new \ReflectionClassConstant(ConstantParentProbe::class, 'PRIVATE_VALUE');
if (!$constant->isPrivate() || $constant->getValue() !== 3
    || $constant->getDeclaringClass()->getName() !== ConstantParentProbe::class) {
    throw new \Exception('Private constant reflection failed');
}
echo "class constants reflection passed\n";
