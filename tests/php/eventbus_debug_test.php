<?php
namespace tests\php;
$listeners = [];
$listeners[] = function ($method) {
    $f = function ($return) use ($method) {
        return "x";
    };
    \Log::info("inner_type=".gettype($f)." callable=".(is_callable($f)?"yes":"no"));
    return $f;
};
$result = $listeners[0]("authenticate");
\Log::info("result_type=".gettype($result)." is_null=".(is_null($result)?"yes":"no")." callable=".(is_callable($result)?"yes":"no"));
if ($result !== null) {
    \Log::info("!==null yes");
} else {
    \Log::info("!==null no");
}
$middlewares = [];
if ($result !== null) {
    $middlewares[] = $result;
}
\Log::info("mw_count=".count($middlewares));
$finish = function ($forward = null) use ($middlewares) {
    \Log::info("finish_mw=".count($middlewares));
    foreach ($middlewares as $i => $finisher) {
        \Log::info("finisher_$i type=".gettype($finisher)." callable=".(is_callable($finisher)?"yes":"no"));
        if ($finisher === null) { \Log::info("skip null"); continue; }
        $r = $finisher($forward);
        \Log::info("finisher_ret=".var_export($r,true));
    }
};
$finish(new \stdClass());
\Log::info("done");
