<?php
class SplIterFam_RecursiveArrayIterator extends \ArrayIterator implements \RecursiveIterator {
    public function hasChildren(): bool {
        $cur = $this->current();
        return is_array($cur) || $cur instanceof \Traversable;
    }
    public function getChildren(): SplIterFam_RecursiveArrayIterator {
        $cur = $this->current();
        if ($cur instanceof \Traversable) {
            return new SplIterFam_RecursiveArrayIterator(iterator_to_array($cur));
        }
        return new SplIterFam_RecursiveArrayIterator($cur);
    }
}

$tree = new SplIterFam_RecursiveArrayIterator([
    'root' => ['child1', 'child2'],
]);
$rti = new \RecursiveTreeIterator($tree);
$rti->rewind();
echo "valid=" . ($rti->valid() ? '1' : '0') . "\n";
echo "count=0\n";
$n = 0;
foreach ($rti as $v) {
    $n++;
    echo "item $n: " . json_encode($v) . "\n";
}
echo "total=$n\n";
