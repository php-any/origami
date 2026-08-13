<?php

namespace tests\php;

/**
 * match 分支：模式与 => 之间夹注释时，不应被 ASI 误插分号。
 */

$v = 2;
$r = match ($v) {
    1 => 'one',
    2
        /** comment */
        => 'two',
    default => 'other',
};
if ($r !== 'two') {
    Log::fatal('match comment between pattern and => failed: ' . var_export($r, true));
}

// ramsey/uuid UuidBuilder 风格：常量表达式 + phpstan 注释 + =>
class MatchCommentArm_UuidLike
{
    public const TYPE_A = 1;
    public const TYPE_B = 2;

    public static function map(int $v): string
    {
        return match ($v) {
            /** @phpstan-ignore possiblyImpure.new */
            self::TYPE_A => 'a',
            self::TYPE_B
                /** @phpstan-ignore possiblyImpure.new */
                => 'b',
            default => 'x',
        };
    }
}

if (MatchCommentArm_UuidLike::map(2) !== 'b') {
    Log::fatal('uuid-like match arm failed');
}

Log::info('match comment arm 测试通过');
