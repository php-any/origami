<?php

$randomizer = new Random\Randomizer();

$keys = $randomizer->pickArrayKeys(['a' => 10, 'b' => 20, 'c' => 30], 2);
if (count($keys) !== 2 || $keys[0] === $keys[1]) {
    Log::fatal('Randomizer::pickArrayKeys 未返回两个不同键');
}
foreach ($keys as $key) {
    if (!in_array($key, ['a', 'b', 'c'], true)) {
        Log::fatal('Randomizer::pickArrayKeys 返回了未知键');
    }
}

$shuffled = $randomizer->shuffleArray([1, 2, 3, 4]);
sort($shuffled);
if ($shuffled !== [1, 2, 3, 4]) {
    Log::fatal('Randomizer::shuffleArray 丢失了元素');
}

$bytes = $randomizer->getBytesFromString('abc', 16);
if (strlen($bytes) !== 16 || preg_match('/[^abc]/', $bytes)) {
    Log::fatal('Randomizer::getBytesFromString 返回值无效');
}

Log::info('randomizer_laravel_test OK');
