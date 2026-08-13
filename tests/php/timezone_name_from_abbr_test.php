<?php

namespace tests\php;

/**
 * timezone_name_from_abbr 函数测试：
 * - 根据缩写返回 IANA 时区名
 * - 根据 UTC 偏移与是否 DST 返回时区名
 * - 无匹配时返回 false
 */

// 1. 按缩写匹配：CET 应返回某欧洲时区（如 Europe/Berlin）
$name = timezone_name_from_abbr('CET');
if ($name === false || $name === '') {
    Log::fatal('timezone_name_from_abbr("CET") 应返回时区名，得到: ' . var_export($name, true));
}
if (!is_string($name)) {
    Log::fatal('timezone_name_from_abbr("CET") 应返回 string，得到: ' . gettype($name));
}

// 2. 按偏移匹配：3600 秒、非 DST（冬季），应返回某 UTC+1 时区
$name2 = timezone_name_from_abbr('', 3600, 0);
if ($name2 === false || $name2 === '') {
    Log::fatal('timezone_name_from_abbr("", 3600, 0) 应返回时区名，得到: ' . var_export($name2, true));
}

// 3. 无匹配时应返回 false
$bad = timezone_name_from_abbr('NOTANABBR123');
if ($bad !== false) {
    Log::fatal('timezone_name_from_abbr("NOTANABBR123") 应返回 false，得到: ' . var_export($bad, true));
}

Log::info('timezone_name_from_abbr 函数测试通过');
