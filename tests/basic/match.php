<?php
namespace tests\basic;

$num = 99;

$ret = match ($num) {
    0 => "zero",
    1 => "one",
    99 => "two",
    default => "many"
};

if ($ret == "two") {
    Log::info("match语句正常匹配");
} else{
    Log::fatal("match语句匹配异常");
}

$ret = match ($num) {
    0 => "zero",
    1 => "one",
    2 => "two",
    default => "many"
};

if ($ret == "many") {
    Log::info("match语句默认值正常匹配");
} else{
    Log::fatal("match语句默认值匹配异常");
}

try {
$ret = match ($num) {
    0 => "zero",
    1 => "one",
    2 => "two"
};
Log::fatal("未匹配的 match 没有抛出异常");
} catch (\UnhandledMatchError $e) {
    Log::info("match 未匹配抛出 UnhandledMatchError");
}
