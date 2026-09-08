<?php

namespace tests\php;

/**
 * is_nan / is_infinite / is_finite，供 brick/math BigDecimal 校验 float。
 */
if (is_nan(\NAN) !== true) {
    Log::fatal('is_nan(NAN) 应为 true');
}
if (is_nan(1.5) !== false) {
    Log::fatal('is_nan(1.5) 应为 false');
}
if (is_infinite(\INF) !== true) {
    Log::fatal('is_infinite(INF) 应为 true');
}
if (is_infinite(1.5) !== false) {
    Log::fatal('is_infinite(1.5) 应为 false');
}
if (is_finite(1.5) !== true) {
    Log::fatal('is_finite(1.5) 应为 true');
}
if (is_finite(\NAN) !== false || is_finite(\INF) !== false) {
    Log::fatal('is_finite(NAN/INF) 应为 false');
}

Log::info('is_nan/is_infinite/is_finite 测试通过');
