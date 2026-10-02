<?php
function captureCheck($condition, $message) {
    if (!$condition) {
        throw new Exception($message);
    }
}
class CaptureAsset {
    public string $package = 'initial';
}
$asset = new CaptureAsset();
$assets = [$asset, $asset];
$counter = 0;
$inner = function () use ($assets, &$counter) {
    $assets[0]->package = 'request';
    $counter++;
    return $assets;
};
$outer = function () use ($inner, &$counter) {
    $result = $inner();
    $counter++;
    return $result;
};
$read = function () use (&$counter) { return $counter; };
$result = $outer();
captureCheck($result[0] === $result[1] && $result[0] === $asset, 'captured object aliases');
captureCheck($asset->package === 'request' && $read() === 2, 'nested captures and shared references');
$keys = [100 => 'removed'];
unset($keys[100]);
$append = function () use ($keys) { $keys[] = 'next'; return $keys; };
captureCheck(array_key_first($append()) === 101 && count($keys) === 0, 'captured array key history');
echo "PASS\n";
