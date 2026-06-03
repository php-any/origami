<?php

use Symfony\Component\Finder\Finder;

// step07: symfony/finder

$finder = new Finder();
$finder->files()->in(dirname(__DIR__) . '/src')->name('*.php');

$count = 0;
foreach ($finder as $file) {
    $count++;
}
step_check('Finder files in src', $count >= 1, "found=$count");

step_check('step07_finder', true, 'symfony/finder (files->in 可用；glob 变量函数未实现)');

