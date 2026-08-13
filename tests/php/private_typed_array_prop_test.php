<?php

namespace tests\php;

class PTA_Holder
{
    private array $dirs = [];

    public function add($d)
    {
        $this->dirs = array_merge($this->dirs, [$d]);
        return $this;
    }

    public function getDirs()
    {
        return $this->dirs;
    }

    public function countDirs()
    {
        return count($this->dirs);
    }
}

$h = new PTA_Holder();
$h->add('/tmp');
if ($h->countDirs() !== 1) {
    Log::fatal('private array 属性写入失败 count=' . $h->countDirs() . ' dirs=' . var_export($h->getDirs(), true));
}
if ($h->getDirs()[0] !== '/tmp') {
    Log::fatal('private array 内容错误');
}
Log::info('private_typed_array_prop 测试通过');
