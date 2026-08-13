<?php

class WcBox { public function bar(): int { return 7; } }
function wc_inner(): WcBox { return new WcBox(); }
function wc_wrap($x): WcBox { return wc_inner(); }

class WcRW {
    public function t(): array {
        return ['value' => wc_wrap(1)->bar()];
    }
    public function deep(): array {
        return ['value' => wc_wrap(wc_inner()->bar())->bar()];
    }
}

$w = new WcRW();
$a = $w->t();
if (($a['value'] ?? 0) !== 7) Log::fatal('t');
$b = $w->deep();
if (($b['value'] ?? 0) !== 7) Log::fatal('deep');
Log::info('wrap chain ok');
