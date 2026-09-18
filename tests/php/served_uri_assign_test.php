<?php

namespace tests\php;

/**
 * FilesystemServiceProvider::serveFiles：方法返回后 booted 回调才执行
 * use (&$served) 的 $served[$uri] = $disk。父帧被 Context 池回收后引用必须仍有效。
 */

class ServedUri_Provider
{
    public $callbacks = [];
    public $result;

    public function serveFiles()
    {
        $served = [];
        $disks = [
            'local' => '/storage',
            'public' => '/public-storage',
        ];
        foreach ($disks as $disk => $uri) {
            $this->callbacks[] = function ($app) use ($disk, $uri, &$served) {
                if (isset($served[$uri])) {
                    throw new \InvalidArgumentException("conflict at {$uri}");
                }
                $served[$uri] = $disk;
                $app->result = $served;
            };
        }
        $this->touch();
    }

    public function touch()
    {
        $a = 1;
        $b = 2;
        return $a + $b;
    }

    public function boot()
    {
        foreach ($this->callbacks as $cb) {
            $cb($this);
        }
    }
}

class ServedUri_Churn
{
    public function go($i)
    {
        $x = $i;
        $y = $i * 2;
        $z = $x + $y;
        return $z;
    }
}

$p = new ServedUri_Provider();
$p->serveFiles();
for ($i = 0; $i < 64; $i++) {
    (new ServedUri_Churn())->go($i);
}
$p->boot();
if (($p->result['/storage'] ?? null) !== 'local') {
    Log::fatal('方法返回后 use (&$served) 丢失 local: '.var_export($p->result, true));
}
if (($p->result['/public-storage'] ?? null) !== 'public') {
    Log::fatal('方法返回后 use (&$served) 未共享 public: '.var_export($p->result, true));
}

$served = [];
$fn = static function () use (&$served) {
    $url = 'http://localhost';
    $uri = rtrim(parse_url($url)['path'] ?? '', '/');
    if ($uri === '') {
        $uri = '/storage';
    }
    $served[$uri] = 'local';
    return $served;
};
$got = $fn();
if (($got['/storage'] ?? null) !== 'local') {
    Log::fatal('use (&$served) 赋值失败: '.var_export($got, true));
}

$served2 = [];
$fn2 = static function () use (&$served2) {
    $uri = rtrim((string) (parse_url('https://example.com')['path'] ?? ''), '/');
    $served2[$uri] = 'public';
    return $served2;
};
$got2 = $fn2();
if (($got2[''] ?? $got2[null] ?? null) !== 'public' && ($got2['/'] ?? null) !== 'public') {
    if (!array_key_exists('', $got2) && !array_key_exists(null, $got2)) {
        Log::fatal('空 path 赋值失败: '.var_export($got2, true));
    }
}

Log::info('served_uri_assign 测试通过');
