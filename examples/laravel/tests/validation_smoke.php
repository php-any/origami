<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Container\Container;
use Illuminate\Translation\ArrayLoader;
use Illuminate\Translation\Translator;
use Illuminate\Validation\Factory;

$container = new Container();
$loader = new ArrayLoader();
$translator = new Translator($loader, 'en');
$factory = new Factory($translator, $container);

$failValidator = $factory->make(['title' => ''], ['title' => 'required|min:1']);
$passValidator = $factory->make(['title' => 'hello'], ['title' => 'required|min:1|max:200']);

if (!$failValidator->fails()) {
    echo "FAIL: 空 title 应校验失败\n";
    exit(1);
}

if (!$passValidator->passes()) {
    echo "FAIL: 合法 title 应校验通过\n";
    foreach ($passValidator->errors()->all() as $msg) {
        echo "  - $msg\n";
    }
    exit(1);
}

echo "PASS\n";
