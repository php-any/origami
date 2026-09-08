<?php

namespace tests\php;

/**
 * DateTime::diff 必须返回真实 DateInterval（y/m/d/h/i/s/invert/days），
 * 且 DateInterval::format / createFromDateString 对齐 PHP，供 Carbon 使用。
 */
$from = new \DateTime('2020-01-01 00:00:00');
$to = new \DateTime('2021-03-05 04:06:08');
$iv = $from->diff($to);

if (!($iv instanceof \DateInterval)) {
    Log::fatal('DateTime::diff 未返回 DateInterval: ' . gettype($iv));
}
if ((int) $iv->y !== 1 || (int) $iv->m !== 2 || (int) $iv->d !== 4) {
    Log::fatal('DateTime::diff 日历分量错误: y=' . $iv->y . ' m=' . $iv->m . ' d=' . $iv->d);
}
if ((int) $iv->h !== 4 || (int) $iv->i !== 6 || (int) $iv->s !== 8) {
    Log::fatal('DateTime::diff 时分秒错误: h=' . $iv->h . ' i=' . $iv->i . ' s=' . $iv->s);
}
if ((int) $iv->invert !== 0) {
    Log::fatal('正向 diff invert 应为 0, got ' . $iv->invert);
}
if ($iv->days !== false && (int) $iv->days < 400) {
    Log::fatal('DateTime::diff days 异常: ' . var_export($iv->days, true));
}

$back = $to->diff($from);
if ((int) $back->invert !== 1) {
    Log::fatal('反向 diff invert 应为 1, got ' . $back->invert);
}
$abs = $to->diff($from, true);
if ((int) $abs->invert !== 0) {
    Log::fatal('absolute diff invert 应为 0, got ' . $abs->invert);
}
if ($from->diff($to)->format('%r%y') !== '1') {
    Log::fatal('DateInterval::format %r%y 失败: ' . $from->diff($to)->format('%r%y'));
}
if ($to->diff($from)->format('%r%y') !== '-1') {
    Log::fatal('负向 DateInterval::format %r%y 失败: ' . $to->diff($from)->format('%r%y'));
}

$micro = \DateInterval::createFromDateString('196942 microseconds');
if (!($micro instanceof \DateInterval)) {
    Log::fatal('createFromDateString(microseconds) 失败');
}
$us = (int) round(((float) $micro->f) * 1000000);
if ($us < 196900 || $us > 196980) {
    Log::fatal('createFromDateString 微秒写入 f 失败: f=' . $micro->f . ' us=' . $us);
}

$day = \DateInterval::createFromDateString('1 day 2 hours');
if ((int) $day->d !== 1 || (int) $day->h !== 2) {
    Log::fatal('createFromDateString(1 day 2 hours) 失败: d=' . $day->d . ' h=' . $day->h);
}

$dt = new \DateTime('2020-01-01');
$dt->add(new \DateInterval('P1D'));
if ($dt->format('Y-m-d') !== '2020-01-02') {
    Log::fatal('DateTime::add(DateInterval) 失败: ' . $dt->format('Y-m-d'));
}

$uu = (new \DateTime('2020-01-01 00:00:00'))->format('U.u');
$fromU = \DateTime::createFromFormat('U.u', $uu);
if ($fromU === false || !($fromU instanceof \DateTime)) {
    Log::fatal('DateTime::createFromFormat(U.u) 失败: ' . var_export($fromU, true));
}
if ($fromU->format('Y-m-d') !== '2020-01-01') {
    Log::fatal('createFromFormat(U.u) 日期错误: ' . $fromU->format('Y-m-d H:i:s') . ' raw=' . $uu);
}

Log::info('DateTime::diff / DateInterval 测试通过');
