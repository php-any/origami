<?php
namespace tests\func;

// PHP uses one array return value and destructuring for several results.
function greet2(): array {
    return ["abc", 123];
}

[$str, $num] = greet2();

if ($num === 123) {
    Log::info("多返回值; 正常");
} else {
    Log::fatal("多返回值; 异常");
}

if ($str === "abc") {
    Log::info("多返回值; 正常");
} else {
    Log::fatal("多返回值; 异常");
}
