<?php
namespace SyntaxMixedGroupSource;
class Item { public function value(): int { return 1; } }
function value(): int { return 2; }
const VALUE = 3;

namespace SyntaxMixedGroupConsumer;
use SyntaxMixedGroupSource\{Item as LocalItem, function value as localValue, const VALUE as LOCAL_VALUE};
if ((new LocalItem)->value() !== 1 || localValue() !== 2 || LOCAL_VALUE !== 3) {
    throw new \Exception('mixed group import kinds');
}
echo "mixed group imports OK\n";
