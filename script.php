namespace tests\func;

function div($obj) {
    return "<div>" + $obj->body + "</div>";
}

function span($obj) {
    return "<span>" + $obj->body + "</span>";
}

$html = div {
    "body": span {
        "body": "内容",
    }
}

if("<div><span>内容</span></div>" == $html) {
    Log::info("函数参数后置; 正常");
} else {
    Log::fatal("函数参数后置; 异常");
}