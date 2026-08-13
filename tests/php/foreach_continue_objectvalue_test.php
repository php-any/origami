<?php

namespace tests\php;

/**
 * foreach 关联数组（ObjectValue）中全部 continue 后，ContinueControl 不得泄漏出循环。
 * 复现 Laravel AgentDetector::fromKnownEnvVars。
 */

enum ForeachContinue_Known: string
{
    case Claude = 'claude';
    case Other = 'other';
}

class ForeachContinue_Detector
{
    public const VARS = [
        'A' => ForeachContinue_Known::Claude,
        'B' => ForeachContinue_Known::Other,
    ];

    public static function allContinue(): string
    {
        foreach (self::VARS as $envVar => $agent) {
            if (getenv($envVar) === false) {
                continue;
            }
            return 'hit';
        }
        return 'none';
    }

    public static function skipFirst(): string
    {
        foreach (self::VARS as $k => $v) {
            if ($k === 'A') {
                continue;
            }
            return $v->value;
        }
        return 'none';
    }

    public static function returnInside(): string
    {
        foreach (self::VARS as $k => $v) {
            return $v->value;
        }
        return 'none';
    }
}

if (ForeachContinue_Detector::allContinue() !== 'none') {
    Log::fatal('全部 continue 应返回 none');
}
if (ForeachContinue_Detector::skipFirst() !== 'other') {
    Log::fatal('跳过首项应得到 other');
}
if (ForeachContinue_Detector::returnInside() !== 'claude') {
    Log::fatal('循环内 return 应得到 claude');
}

$r = ForeachContinue_Detector::allContinue();
$r2 = match (true) {
    $r === 'none' => 'ok',
    default => 'bad',
};
if ($r2 !== 'ok') {
    Log::fatal('continue 后 match 失败');
}

Log::info('foreach_continue_objectvalue 测试通过');
