<?php

class McBox
{
    public function foo(): McBox
    {
        return $this;
    }

    public function bar(): int
    {
        return 42;
    }
}

function mc_maker(): McBox
{
    return new McBox();
}

class McRW
{
    public function chain(): array
    {
        return ['value' => mc_maker()->bar()];
    }

    public function nested(): array
    {
        return ['value' => mc_maker()->foo()->bar()];
    }
}

$w = new McRW();
$a = $w->chain();
if (($a['value'] ?? 0) !== 42) {
    Log::fatal('chain fail');
}
$b = $w->nested();
if (($b['value'] ?? 0) !== 42) {
    Log::fatal('nested fail');
}
Log::info('array literal method chain ok');
