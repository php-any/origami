<?php
/**
 * 对照 PHP 8.4 黄金基准（storage/origami-debug/php84/php.exe）逐行 diff，
 * 覆盖 Collection / Macroable 的补齐项：mixin、实例 collect()、proxy()、
 * ensure()、getCachingIterator()，以及「方法/属性不存在」的异常类名。
 *
 * 已知且可接受的差异（只有这两行）：
 *   dump  —— PHP 走 Symfony VarDumper（带 ANSI 着色与文件行号），origami 是纯文本；
 *   collect_static —— PHP 对「静态调用非静态方法」抛 Error，origami 的
 *     GetStaticMethod 目前会把任意实例方法当作可静态调用（跨全部桩类的既有行为，
 *     不在本次改动范围）。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use Illuminate\Support\Collection;

function t($l, $fn)
{
    echo $l.' => ';
    try { echo json_encode($fn()); } catch (Throwable $e) { echo 'EX '.get_class($e).' '.$e->getMessage(); }
    echo "\n";
}

t('sortBy', fn () => collect([3, 1, 2])->sortBy(fn ($v) => $v)->values()->all());
t('sortByDesc', fn () => collect([3, 1, 2])->sortByDesc(fn ($v) => $v)->values()->all());
t('sortByMany', fn () => collect([['a' => 2, 'b' => 1], ['a' => 1, 'b' => 2]])->sortBy([['a', 'asc']])->values()->all());
t('dump', fn () => collect([1])->dump() instanceof Collection);
t('dd_exists', fn () => method_exists(Collection::class, 'dd'));
t('collect_static', fn () => Collection::collect([1, 2])->all());
t('mixin', fn () => method_exists(Collection::class, 'mixin'));
t('lazy_only', fn () => method_exists(Collection::class, 'lazy'));
t('sliding', fn () => collect([1, 2, 3, 4])->sliding(2)->map(fn ($w) => $w->all())->all());
t('chunkWhile', fn () => collect([1, 2, 3])->chunkWhile(fn ($v, $k, $c) => $v > 1)->map(fn ($w) => $w->all())->all());
t('chunkWhile_assoc', fn () => collect(['a' => 1, 'b' => 2, 'c' => 3])->chunkWhile(fn ($v, $k, $c) => $v > 1)->map(fn ($w) => $w->all())->all());
t('chunkWhile_chunk_arg', fn () => collect([1, 2, 3])->chunkWhile(fn ($v, $k, $c) => $c->count() < 2)->map(fn ($w) => $w->all())->all());
t('flatten_depth', fn () => collect([1, [2, [3]]])->flatten()->all());
t('median', fn () => collect([1, 2, 3])->median());
t('whereInstanceOf', fn () => collect([1, new stdClass])->whereInstanceOf(stdClass::class)->count());
t('tap_each', function () {
    $c = collect([1, 2]);
    return $c->tap(fn ($x) => null)->all();
});
t('toHtml', fn () => collect(['<b>'])->toHtml());
t('times_static', fn () => Collection::times(3, fn ($i) => $i)->all());
t('range_static', fn () => Collection::range(1, 3)->all());
t('make_static', fn () => Collection::make([1])->all());
t('wrap_static', fn () => Collection::wrap(5)->all());
t('unwrap_static', fn () => Collection::unwrap(collect([1])));

// ---- mixin / 实例 collect / proxy / ensure / getCachingIterator / 异常类名 ----

class _MixinProbe
{
    public function mirrored() { return function ($v) { return $v; }; }
    public function shout() { return function ($s) { return strtoupper($s).'!'; }; }
    public function dup() { return function () { return 'from-mixin'; }; }
    protected function prot() { return function () { return 'protected'; }; }
    private function hidden() { return function () { return 'hidden'; }; }
}

class _MixinChildProbe extends _MixinProbe
{
    public function extra() { return function () { return 'extra'; }; }
}

class _StrMixinProbe
{
    public function excited() { return function ($s) { return strtoupper($s).'!'; }; }
}

class _ReqMixinProbe
{
    public function probeHello() { return function () { return 'hello-request'; }; }
}

class _EnsureBase {}

class _EnsureChild extends _EnsureBase {}

t('mixin_register', function () {
    Collection::mixin(new _MixinProbe);
    return [Collection::hasMacro('mirrored'), Collection::hasMacro('hidden')];
});
t('mixin_call_instance', fn () => collect([1])->mirrored(7));
t('mixin_call_static', fn () => Collection::shout('hi'));

t('mixin_replace_true', function () {
    Collection::macro('dup', fn () => 'pre-existing');
    Collection::mixin(new _MixinProbe, true);
    return collect([1])->dup();
});
t('mixin_replace_false', function () {
    Collection::macro('dup', fn () => 'pre-existing');
    Collection::mixin(new _MixinProbe, false);
    return collect([1])->dup();
});
t('mixin_replace_false_new', function () {
    Collection::flushMacros();
    Collection::mixin(new _MixinProbe, false);
    return collect([1])->dup();
});
t('mixin_inherited', function () {
    Collection::flushMacros();
    Collection::mixin(new _MixinChildProbe);
    return [
        Collection::hasMacro('extra'),
        Collection::hasMacro('mirrored'),
        Collection::hasMacro('prot'),
        Collection::hasMacro('hidden'),
    ];
});

t('collect_instance', fn () => collect([1, 2])->collect()->all());
t('collect_instance_class', fn () => get_class(collect([1, 2])->collect()));

t('magicget_unknown', function () {
    try { return collect([1])->countBy; } catch (Throwable $e) { return [get_class($e), $e->getMessage()]; }
});
t('proxy_register', function () {
    Collection::proxy('countBy');
    return get_class(collect([1])->countBy);
});
t('proxy_call', fn () => collect([['g' => 'a'], ['g' => 'b'], ['g' => 'a']])->countBy->g->all());

t('ensure_pass', fn () => collect([1, 2])->ensure('int')->all());
t('ensure_same_instance', fn () => (function () { $c = collect([1]); return $c->ensure('int') === $c; })());
t('ensure_fail', function () {
    try { return collect([1, 'a'])->ensure('int'); } catch (Throwable $e) { return [get_class($e), $e->getMessage()]; }
});
t('ensure_types_array', fn () => collect([1, 'a'])->ensure(['int', 'string'])->count());
t('ensure_types_array_fail', function () {
    try { return collect([1, 2.5])->ensure(['int', 'string']); }
    catch (Throwable $e) { return [get_class($e), $e->getMessage()]; }
});
t('ensure_inherited', function () {
    Collection::proxy('pad');   // 顺带证明 proxy 可重复注册
    return collect([new _EnsureChild])->ensure(_EnsureBase::class)->count();
});
t('ensure_fail_inherited', function () {
    try { return collect([new _EnsureChild, 1])->ensure(_EnsureBase::class); }
    catch (Throwable $e) { return [get_class($e), $e->getMessage()]; }
});

t('caching_iterator_class', fn () => get_class(collect([1, 2])->getCachingIterator()));
t('caching_iterator_index', function () {
    $it = collect(['a', 'b'])->getCachingIterator();
    $it->rewind();
    return [(string) $it, $it->hasNext()];
});

t('badmethodcall_class', function () {
    try { return collect([1])->noSuchMethodXyz(); } catch (Throwable $e) { return [get_class($e), $e->getMessage()]; }
});
t('badmethodcall_is_logic', function () {
    try { collect([1])->noSuchMethodXyz(); return 'no-throw'; } catch (LogicException $e) { return 'caught'; }
});
t('badmethodcall_static', function () {
    try { return Collection::noSuchStaticXyz(); } catch (Throwable $e) { return [get_class($e), $e->getMessage()]; }
});

t('str_mixin', function () {
    Illuminate\Support\Str::mixin(new _StrMixinProbe);
    return [Illuminate\Support\Str::hasMacro('excited'), Illuminate\Support\Str::excited('ok')];
});

t('request_mixin', function () {
    Illuminate\Http\Request::mixin(new _ReqMixinProbe);
    return [Illuminate\Http\Request::hasMacro('probeHello'), Illuminate\Http\Request::probeHello()];
});

echo "DONE\n";
