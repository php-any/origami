<?php
namespace SyntaxGroupClassSource;
class First { public function value(): int { return 1; } }
class Second { public function value(): int { return 2; } }

namespace SyntaxGroupClassConsumer;
use SyntaxGroupClassSource\{First, Second as Renamed,};
if ((new First)->value() !== 1 || (new Renamed)->value() !== 2
    || First::class !== 'SyntaxGroupClassSource\First'
    || Renamed::class !== 'SyntaxGroupClassSource\Second') {
    throw new \Exception('group class imports and alias');
}
echo "group class imports OK\n";
