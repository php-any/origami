<?php

namespace tests\hash;

/**
 * hash：md5 / sha1 与 hash() 对照。
 */

if (md5('x') !== hash('md5', 'x')) {
    Log::fatal('md5 与 hash(md5) 不一致');
}
if (sha1('x') !== hash('sha1', 'x')) {
    Log::fatal('sha1 与 hash(sha1) 不一致');
}

Log::info('hash md5/sha1 测试通过');
