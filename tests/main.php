<?php

namespace App;

class App
{
    public function main()
    {
        echo "OK";
    }
}

function div(config) {
    echo config, "\n"


    config.body(123)
}

div {
    width: 100px,
    "body": (a) => {
        echo "body" + a;
    }
}

echo dump(time(), 123);

用户 = "张三"

输出 用户

函数 乘法(左，右) {
    输出 左 × 右
}

乘法(3, 5)