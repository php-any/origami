<?php
namespace tests\php;
$fn = function ($event) { return $event; };
$wrapped = $fn(...);
if (!is_callable($wrapped)) Log::fatal('not callable');
if ($wrapped('ok') !== 'ok') Log::fatal('invoke fail');
Log::info('first_class_closure_callable 测试通过');
