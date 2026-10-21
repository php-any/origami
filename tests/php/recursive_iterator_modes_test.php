<?php
function traverseTree($mode, $depth = -1) {
    $it = new RecursiveIteratorIterator(new RecursiveArrayIterator(['a' => ['b' => 1, 'c' => ['d' => 2]], 'e' => 3]), $mode);
    $it->setMaxDepth($depth);
    $result = [];
    foreach ($it as $key => $value) { if ($it->getInnerIterator()->key() !== $key) { throw new Exception('Incorrect inner iterator'); } $result[] = $key . ':' . $it->getDepth(); }
    return implode(',', $result);
}
foreach ([
    [0, -1, 'b:1,d:2,e:0'],
    [1, -1, 'a:0,b:1,c:1,d:2,e:0'],
    [2, -1, 'b:1,d:2,c:1,a:0,e:0'],
    [0, 0, 'e:0'],
    [1, 0, 'a:0,e:0'],
    [2, 0, 'a:0,e:0'],
    [0, 1, 'b:1,e:0']
] as $case) {
    $actual = traverseTree($case[0], $case[1]);
    if ($actual !== $case[2]) { throw new Exception('Traversal mismatch: ' . $actual); }
}
echo "recursive iterator modes passed\n";
