<?php

namespace tests\php;

$tz = timezone_open('Europe/Paris');

// 用 serialize 看看内部属性
Log::info('serialized tz = ' . serialize($tz));
