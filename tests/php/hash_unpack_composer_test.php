<?php
foreach (['' => '2d06800538d394c2', '1' => '65cd25028f98f158', 'hello' => '9555e8555c62dcfd'] as $input => $expected) {
    if (hash('xxh3', (string) $input) !== $expected || bin2hex(hash('xxh3', (string) $input, true)) !== $expected) {
        throw new Exception('XXH3 digest mismatch');
    }
}
if (hash('xxh3', str_repeat('abc', 1000)) !== '1fa2f942314cd628') { throw new Exception('Long XXH3 mismatch'); }
try { hash('not-an-algorithm', 'data'); throw new Exception('Unknown hash accepted'); } catch (ValueError $e) {}
if (unpack('ihash', "\xff\xff\xff\xff") !== ['hash' => -1]) { throw new Exception('Signed unpack failed'); }
if (unpack('n2value/Ctail', "x\x00\x01\x01\x00\x07", 1) !== ['value1' => 1, 'value2' => 256, 'tail' => 7]) { throw new Exception('Offset unpack failed'); }
if (unpack('C*byte', "\x01\x02") !== ['byte1' => 1, 'byte2' => 2]) { throw new Exception('Repeated unpack failed'); }
$a = ['z' => 3, 'a' => 1, 'm' => 2]; $copy = $a;
uksort($a, function ($a, $b) { return strcmp($a, $b); });
if (array_keys($a) !== ['a', 'm', 'z'] || array_keys($copy) !== ['z', 'a', 'm']) { throw new Exception('Key sort or copy failed'); }
echo "hash/unpack/key sort passed\n";
