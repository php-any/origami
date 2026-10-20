<?php

namespace tests\php;

/**
 * 类同时 use 子 trait 和它已经 use 过的 trait 时，同一个方法不是冲突。
 */

trait TraitNestedReuse_Files {
    protected function writeFile(): string {
        return 'file';
    }
}

trait TraitNestedReuse_Panels {
    use TraitNestedReuse_Files;

    public function panel(): string {
        return $this->writeFile();
    }
}

class TraitNestedReuse_Command {
    use TraitNestedReuse_Panels;
    use TraitNestedReuse_Files;

    public function run(): string {
        return $this->panel();
    }
}

$got = (new TraitNestedReuse_Command())->run();
if ($got !== 'file') {
    \Log::fatal('嵌套 trait 再次 use 应复用同一方法，实际 ' . $got);
}

\Log::info('trait nested reuse 测试通过');
