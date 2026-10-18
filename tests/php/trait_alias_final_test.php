<?php
trait FinalAliasTrait { final public function value(): int { return 7; } }
class FinalAliasTarget {
    use FinalAliasTrait;
}
if ((new FinalAliasTarget)->value() !== 7
    || !(new ReflectionMethod(FinalAliasTarget::class, 'value'))->isFinal()) {
    throw new Exception('final trait method must stay final after use');
}
echo "trait alias final OK\n";
