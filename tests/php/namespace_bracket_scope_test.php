<?php
namespace SyntaxBracketOne {
    class Item { public function name(): string { return 'one'; } }
    function name(): string { return 'one'; }
    const LABEL = 'one';
}
namespace SyntaxBracketTwo {
    use SyntaxBracketOne\Item as ImportedItem;
    class Item { public function name(): string { return 'two'; } }
    function name(): string { return 'two'; }
    const LABEL = 'two';
    if ((new ImportedItem)->name() !== 'one' || (new Item)->name() !== 'two'
        || name() !== 'two' || LABEL !== 'two') {
        throw new \Exception('bracketed namespace lookup or alias scope');
    }
}
namespace {
    if ((new \SyntaxBracketOne\Item)->name() !== 'one'
        || \SyntaxBracketOne\name() !== 'one' || \SyntaxBracketOne\LABEL !== 'one'
        || \SyntaxBracketTwo\name() !== 'two') {
        throw new \Exception('global namespace block lookup');
    }
    echo "bracketed namespaces OK\n";
}
