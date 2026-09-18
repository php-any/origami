<?php
namespace tests\php;
$scopes = ['a'];
$fn = function () use ($scopes) {
    return $scopes[0];
};
$got = $fn();
if ($got !== 'a') {
    Log::fatal('closure use failed: '.var_export($got, true));
}
Log::info('closure use ok');
