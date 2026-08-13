<?php

namespace tests\php;

require dirname(__DIR__, 2) . '/examples/laravel/vendor/autoload.php';

use Symfony\Component\Console\Input\ArrayInput;
use Symfony\Component\Console\Output\BufferedOutput;
use Symfony\Component\Console\Style\SymfonyStyle;

$buf = new BufferedOutput();
$style = new SymfonyStyle(new ArrayInput([]), $buf);
$style->success('Formatted success message');
$out = $buf->fetch();
Log::info('raw: ' . var_export($out, true));
Log::info('hex: ' . bin2hex($out));
foreach (str_split($out) as $i => $ch) {
    $o = ord($ch);
    if ($o < 32 || $o === 92) {
        Log::info("pos $i ord=$o");
    }
}
Log::info('symfony_success_block_test 探测完成');
