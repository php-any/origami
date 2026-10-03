<?php
function typedReturnGenerator(bool $stop): Traversable {
    if ($stop) { return; }
    yield 'key' => 'value';
    return 3;
}
if (iterator_to_array(typedReturnGenerator(true)) !== []) { throw new Exception('generator bare return'); }
if (iterator_to_array(typedReturnGenerator(false)) !== ['key' => 'value']) { throw new Exception('generator return declaration'); }
$g = typedReturnGenerator(false);
$caught = false;
try { $g->getReturn(); } catch (Exception $e) { $caught = true; }
if (!$caught) { throw new Exception('generator premature return'); }
iterator_to_array($g);
if ($g->getReturn() !== 3 || $g->valid()) { throw new Exception('generator saved return'); }
$g->next();
if ($g->valid() || $g->getReturn() !== 3) { throw new Exception('generator closed state'); }
echo "typed generator OK\n";
