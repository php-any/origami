<?php

namespace tests\php;

/**
 * 实例方法内调用 class_exists 必须能拿到 VM（方法帧 inner 不得丢掉钉住的 vm）。
 */

class ClassExistsFromMethod_Box
{
    public function check()
    {
        return class_exists(self::class);
    }
}

$box = new ClassExistsFromMethod_Box();
if ($box->check() !== true) {
    Log::fatal('实例方法内 class_exists 失败');
}

Log::info('class_exists 从实例方法测试通过');
