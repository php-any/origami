<?php
class NestedCache {
    private array $cache = [];
    public int $calls = 0;
    public function get(string $key): string {
        return $this->cache[false][$key]['state'] ??= $this->compute();
    }
    private function compute(): string { $this->calls++; return 'cached'; }
    public function snapshot(): array { return $this->cache; }
}
$cache = new NestedCache();
if ($cache->get('40666883792896') !== 'cached'
    || $cache->get('40666883792896') !== 'cached'
    || $cache->calls !== 1) {
    throw new Exception('Nested cache assignment failed');
}
if ($cache->snapshot() !== [0 => ['40666883792896' => ['state' => 'cached']]]) {
    throw new Exception('Boolean cache key must become integer zero');
}
foreach ([false, true, 2.0, null, '08', '12'] as $key) {
    $array = [];
    $array[$key]['middle']['leaf'] = 'one';
    $snapshot = $array;
    $array[$key]['middle']['leaf'] = 'two';
    $expected = [];
    $expected[$key] = ['middle' => ['leaf' => 'two']];
    if ($array !== $expected || $snapshot[$key]['middle']['leaf'] !== 'one') {
        throw new Exception('Recursive writeback key conversion or copy-on-write failed');
    }
}
echo "nested cache coalesce assignment OK\n";
