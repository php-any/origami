<?php
function literalOptions($name) {
    return ['auth' => ['guards', 'providers', 'passwords'], 'cache' => ['stores']][$name] ?? [];
}
if (literalOptions('auth') !== ['guards', 'providers', 'passwords']) {
    throw new Exception('array literal dereference must select the requested element');
}
if (literalOptions('missing') !== []) { throw new Exception('missing literal key must coalesce'); }
if (['outer' => ['inner' => 42]]['outer']['inner'] !== 42) { throw new Exception('nested literal dereference'); }
if ([11, 22][1] !== 22) { throw new Exception('list literal dereference'); }
if (([][0] ?? 'fallback') !== 'fallback') { throw new Exception('empty literal dereference'); }
if ((['value' => null]['value'] ?? 7) !== 7) { throw new Exception('null literal value'); }
$seen = [];
$config = ['guards' => [], 'providers' => [], 'passwords' => []];
foreach (literalOptions('auth') as $option) {
    if (isset($config[$option])) { $seen[] = $option; }
}
if ($seen !== ['guards', 'providers', 'passwords']) { throw new Exception('configuration options iteration'); }
echo "array literal dereference: PASS\n";
