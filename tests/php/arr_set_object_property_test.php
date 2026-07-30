<?php
namespace tests\php;

class ArrSet_Holder
{
    public $items = [];

    public function setNested($key, $value)
    {
        // simulate Arr::set($this->items, ...)
        self::arrSet($this->items, $key, $value);
    }

    public static function arrSet(&$array, $key, $value)
    {
        $keys = explode('.', $key);
        foreach ($keys as $i => $k) {
            if (count($keys) === 1) {
                break;
            }
            unset($keys[$i]);
            if (! isset($array[$k]) || ! is_array($array[$k])) {
                $array[$k] = [];
            }
            $array = &$array[$k];
        }
        $array[array_shift($keys)] = $value;
        return $array;
    }
}

$h = new ArrSet_Holder();
$h->items = ['app' => ['name' => 'Laravel']];
$h->setNested('app.providers', ['A', 'B']);

if (($h->items['app']['providers'] ?? null) !== ['A', 'B']) {
    Log::fatal('对象属性 by-ref Arr::set 失败: '.var_export($h->items, true));
}

Log::info('arr_set_object_property 测试通过');
