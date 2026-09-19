<?php

namespace tests\php;

/**
 * Laravel data_set($object, 'key', $value) 必须写对象属性，不能把对象拷成数组后丢掉。
 */

class DataSetObject_Page
{
    public array $data = [];
}

if (!function_exists('data_set')) {
    \Log::info('data_set 未注册（非 Laravel helper 路径），跳过对象测试');
} else {
    $page = new DataSetObject_Page();
    data_set($page, 'data', ['site_name' => 'Origami Shop', 'order_prefix' => 'ORD']);
    if (($page->data['site_name'] ?? null) !== 'Origami Shop') {
        \Log::fatal('data_set 对象属性未写入: ' . var_export($page->data, true));
    }
    data_set($page, 'data.order_prefix', 'X');
    if (($page->data['order_prefix'] ?? null) !== 'X') {
        \Log::fatal('data_set 点路径未写入: ' . var_export($page->data, true));
    }
}

\Log::info('data_set_object_property 测试通过');
