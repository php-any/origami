<?php
class LiteralConfigProbe {
    public static $config = [
        'bin' => '{$vendor-dir}/bin',
        'cache' => '{$home}/cache',
        'raw' => '$value ${value} {$value}',
    ];
}
if (LiteralConfigProbe::$config['bin'] !== '{' . '$vendor-dir}/bin'
    || LiteralConfigProbe::$config['cache'] !== '{' . '$home}/cache'
    || LiteralConfigProbe::$config['raw'] !== '$value ${value} ' . '{' . '$value}') {
    throw new \Exception('Single quoted strings interpolated');
}
echo "single quoted placeholders passed\n";
