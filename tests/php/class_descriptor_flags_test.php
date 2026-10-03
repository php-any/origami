<?php
trait DescriptorFlagTrait { public function value(): int { return 1; } }
abstract class DescriptorFlagAbstract {}
final class DescriptorFlagFinal {}
readonly class DescriptorFlagReadonly { public int $value; }
enum DescriptorFlagEnum: string { case One = 'one'; }
$expect = [DescriptorFlagAbstract::class => 64, DescriptorFlagFinal::class => 32, DescriptorFlagReadonly::class => 65536, DescriptorFlagEnum::class => 32, DescriptorFlagTrait::class => 0];
foreach ($expect as $class => $modifiers) { if ((new ReflectionClass($class))->getModifiers() !== $modifiers) { throw new Exception('class modifiers lost'); } }
if (!(new ReflectionClass(DescriptorFlagTrait::class))->isTrait() || !(new ReflectionClass(DescriptorFlagEnum::class))->isEnum() || (new ReflectionClass(BackedEnum::class))->isEnum()) { throw new Exception('class category lost'); }
if (!(new ReflectionClass(DescriptorFlagReadonly::class))->isReadOnly() || !(new ReflectionClass(DescriptorFlagFinal::class))->isFinal()) { throw new Exception('class flags lost'); }
echo "class descriptor flags OK\n";
