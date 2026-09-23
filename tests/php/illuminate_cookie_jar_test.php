<?php

if (!class_exists(\Illuminate\Cookie\CookieJar::class, false)) {
    echo "skip: CookieJar 原生类未注册\n";
    return;
}

$jar = new \Illuminate\Cookie\CookieJar();
$jar->setDefaultPathAndDomain('/app', 'example.test', true, 'lax');
$jar->queue('a', '1', 60);
if (!$jar->hasQueued('a')) {
    echo "FAIL hasQueued\n";
    exit(1);
}
$c = $jar->queued('a');
if (!is_object($c)) {
    echo "FAIL queued\n";
    exit(1);
}
$jar->expire('b');
if (!$jar->hasQueued('b')) {
    echo "FAIL expire queues forget cookie\n";
    exit(1);
}
$jar->unqueue('a');
if ($jar->hasQueued('a')) {
    echo "FAIL unqueue\n";
    exit(1);
}

echo "OK illuminate_cookie_jar_test\n";
