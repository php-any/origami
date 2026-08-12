<?php

namespace tests\php;

/**
 * instanceof 右侧 bare 类名须展开 use 别名（复现 tests/obj/B.php）。
 */

use tests\php\instanceof_use_alias_fixtures\InstanceofUseAlias_Parent;

class InstanceofUseAlias_Child extends InstanceofUseAlias_Parent
{
}

$obj = new InstanceofUseAlias_Child();

if (!($obj instanceof InstanceofUseAlias_Parent)) {
    Log::fatal('use 别名 instanceof 失败');
}
if (!($obj instanceof \tests\php\instanceof_use_alias_fixtures\InstanceofUseAlias_Parent)) {
    Log::fatal('FQCN instanceof 失败');
}

Log::info('instanceof_use_alias 测试通过');
