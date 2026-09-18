<?php

namespace tests\php;

/**
 * Symfony ArgvInput::getFirstArgument 遍历 typed private array $tokens。
 */

class ArgvTokens_First
{
    /** @var list<string> */
    private array $tokens;

    public function __construct(array $tokens)
    {
        $this->tokens = $tokens;
    }

    public function getFirstArgument(): ?string
    {
        $isOption = false;
        foreach ($this->tokens as $i => $token) {
            if ($token && '-' === $token[0]) {
                if (str_contains($token, '=') || !isset($this->tokens[$i + 1])) {
                    continue;
                }
                $isOption = true;
                continue;
            }
            if ($isOption) {
                $isOption = false;
                continue;
            }
            return $token;
        }
        return null;
    }
}

$t = new ArgvTokens_First(['inspire']);
$got = $t->getFirstArgument();
if ($got !== 'inspire') {
    Log::fatal('typed array tokens getFirstArgument 失败: '.var_export($got, true));
}

$t2 = new ArgvTokens_First(['serve', '--host=127.0.0.1']);
$got2 = $t2->getFirstArgument();
if ($got2 !== 'serve') {
    Log::fatal('带选项时 first argument 失败: '.var_export($got2, true));
}

Log::info('argv_tokens_first 测试通过');
