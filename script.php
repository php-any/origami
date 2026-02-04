<?php
namespace tests\php;

$nonEmptyString = "hello";
if(!empty($nonEmptyString)) {
    Log::info("非空字符串测试通过");
} else {
    Log::fatal("非空字符串测试失败");
}

