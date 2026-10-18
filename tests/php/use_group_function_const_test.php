<?php
namespace SyntaxGroupValueSource;
function first(): int { return 11; }
function second(): int { return 22; }
const FIRST = 33;
const SECOND = 44;

namespace SyntaxGroupValueConsumer;
use function SyntaxGroupValueSource\{first, second as renamed,};
use const SyntaxGroupValueSource\{FIRST, SECOND as RENAMED,};
if (first() !== 11 || renamed() !== 22 || FIRST !== 33 || RENAMED !== 44) {
    throw new \Exception('group function and constant imports');
}
echo "group function and constant imports OK\n";
