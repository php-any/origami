<?php
namespace tests\php;
if (!preg_match('/(?<attributes>abc)/', 'abc', $m)) {
  Log::fatal('match fail');
}
if (($m['attributes'] ?? '') !== 'abc') {
  Log::fatal('named group missing: '.json_encode($m));
}
Log::info('named group ok');
