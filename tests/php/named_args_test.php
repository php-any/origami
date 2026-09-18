<?php
namespace tests\php;

function greet(string $name, string $greeting = 'hi') {
    return "$greeting $name";
}

$got = greet(greeting: 'hello', name: 'world');
if ($got !== 'hello world') {
    Log::fatal('named args failed: '.$got);
}
Log::info('named args test passed: '.$got);
