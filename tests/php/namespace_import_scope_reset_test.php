<?php
namespace SyntaxImportScopeSource;
class Item {}

namespace SyntaxImportScopeFirst;
use SyntaxImportScopeSource\Item as Alias;
if (Alias::class !== 'SyntaxImportScopeSource\Item') {
    throw new \Exception('first namespace alias');
}

namespace SyntaxImportScopeSecond;
class Alias {}
if (Alias::class !== 'SyntaxImportScopeSecond\Alias'
    || get_class(new Alias) !== 'SyntaxImportScopeSecond\Alias') {
    throw new \Exception('imports must not leak across namespace declarations');
}
echo "namespace import scope reset OK\n";
