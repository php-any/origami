<?php
namespace SyntaxRelativeNames;
const MARKER = 'local';
function marker(): string { return 'function'; }
class Item { public function marker(): string { return 'class'; } }
if (namespace\MARKER !== 'local' || namespace\marker() !== 'function'
    || (new namespace\Item)->marker() !== 'class'
    || namespace\Item::class !== 'SyntaxRelativeNames\Item') {
    throw new \Exception('namespace-relative constant, function and class lookup');
}
echo "namespace relative names OK\n";
