<?php

namespace tests\php;

/**
 * 对齐 PHP：子类实例可传入声明为父类的参数（Laravel LoadConfiguration::getNestedDirectory(SplFileInfo $file)）。
 */
class SplFileInfoHint_Child extends \SplFileInfo
{
}

function SplFileInfoHint_accept(\SplFileInfo $file): string
{
    return $file->getFilename();
}

$child = new SplFileInfoHint_Child(__FILE__);
$got = SplFileInfoHint_accept($child);
if ($got !== basename(__FILE__)) {
    Log::fatal('子类传入 SplFileInfo 参数失败: '.$got);
}

if (!($child instanceof \SplFileInfo)) {
    Log::fatal('子类 instanceof SplFileInfo 应为 true');
}

Log::info('splfileinfo_type_hint 测试通过');
