<?php

namespace tests\origami;

/**
 * vendor 加速层：Illuminate\Support\Arr / Collection 必须在无 include vendor 类文件时可用。
 */

if (!class_exists(\Illuminate\Support\Arr::class, false)) {
    \Log::fatal('Illuminate\\Support\\Arr 应已由 std/illuminate 预注册');
}
if (!class_exists(\Symfony\Component\HttpFoundation\Request::class, false)) {
    \Log::fatal('Symfony Request 应已由 std/symfony/http-foundation 预注册');
}
if (!class_exists(\Symfony\Component\HttpFoundation\Cookie::class, false)) {
    \Log::fatal('Symfony Cookie 应已由 std/symfony/http-foundation 预注册');
}
if (!class_exists(\Symfony\Component\HttpFoundation\JsonResponse::class, false)) {
    \Log::fatal('Symfony JsonResponse 应已由 std/symfony/http-foundation 预注册');
}
if (!class_exists(\Symfony\Component\Finder\Finder::class, false)) {
    \Log::fatal('Symfony Finder 应已由 std/symfony/finder 预注册');
}
if (!class_exists(\Symfony\Component\Routing\Route::class, false)) {
    \Log::fatal('Symfony Route 应已由 std/symfony/routing 预注册');
}
if (!class_exists(\Symfony\Component\Routing\RouteCollection::class, false)) {
    \Log::fatal('Symfony RouteCollection 应已由 std/symfony/routing 预注册');
}
if (!class_exists(\Symfony\Component\String\UnicodeString::class, false)) {
    \Log::fatal('Symfony UnicodeString 应已由 std/symfony/string 预注册');
}
if (!class_exists(\Symfony\Component\String\ByteString::class, false)) {
    \Log::fatal('Symfony ByteString 应已由 std/symfony/string 预注册');
}
if (!class_exists(\Symfony\Component\HttpFoundation\File\UploadedFile::class, false)) {
    \Log::fatal('Symfony UploadedFile 应已由 std/symfony/http-foundation 预注册');
}
if (!class_exists(\Symfony\Component\HttpFoundation\RequestStack::class, false)) {
    \Log::fatal('Symfony RequestStack 应已由 std/symfony/http-foundation 预注册');
}
if (!class_exists(\Symfony\Component\Uid\Uuid::class, false)) {
    \Log::fatal('Symfony Uuid 应已由 std/symfony/uid 预注册');
}
if (!class_exists(\Symfony\Component\Uid\Ulid::class, false)) {
    \Log::fatal('Symfony Ulid 应已由 std/symfony/uid 预注册');
}
if (!class_exists(\Symfony\Component\Clock\NativeClock::class, false)) {
    \Log::fatal('Symfony NativeClock 应已由 std/symfony/clock 预注册');
}
if (!class_exists(\Illuminate\Http\Response::class, false)) {
    \Log::fatal('Illuminate Response 应已由 std/illuminate/http 预注册');
}
if (!function_exists('collect')) {
    \Log::fatal('collect() 应已由 laravel13 vendoraccel 注册');
}

$cookie = \Symfony\Component\HttpFoundation\Cookie::create('k', 'v');
if ($cookie->getName() !== 'k' || $cookie->getValue() !== 'v') {
    \Log::fatal('Cookie::create 语义失败');
}

$jr = new \Symfony\Component\HttpFoundation\JsonResponse(['ok' => true], 201);
if ($jr->getStatusCode() !== 201) {
    \Log::fatal('JsonResponse status 失败');
}
$content = $jr->getContent();
if (!is_string($content) || strpos($content, '"ok"') === false) {
    \Log::fatal('JsonResponse content 失败: '.$content);
}

$rr = new \Symfony\Component\HttpFoundation\RedirectResponse('/admin', 302);
if ($rr->getTargetUrl() !== '/admin') {
    \Log::fatal('RedirectResponse target 失败');
}

$finder = \Symfony\Component\Finder\Finder::create()->files()->name('vendoraccel_smoke.php')->in(__DIR__)->depth(0);
$count = 0;
foreach ($finder as $file) {
    $count++;
    if (!method_exists($file, 'getRealPath') || $file->getFilename() !== 'vendoraccel_smoke.php') {
        \Log::fatal('Finder SplFileInfo 失败');
    }
}
if ($count < 1) {
    \Log::fatal('Finder 未找到本文件');
}

$items = ['x' => 1];
\Illuminate\Support\Arr::set($items, 'y', 2);
if (($items['y'] ?? null) !== 2) {
    \Log::fatal('Arr::set 引用语义失败');
}

$collapsed = \Illuminate\Support\Arr::collapse([['Illuminate\\A', 'Illuminate\\B'], ['Pkg\\C'], ['App\\D']]);
if (!is_array($collapsed) || array_values($collapsed) !== ['Illuminate\\A', 'Illuminate\\B', 'Pkg\\C', 'App\\D']) {
    \Log::fatal('Arr::collapse 应拍平一层数组: '.var_export($collapsed, true));
}

$mapped = \Illuminate\Support\Arr::mapWithKeys(['a' => 1, 'b' => 2], fn ($v, $k) => [$k => $v * 10]);
if (($mapped['a'] ?? null) !== 10 || ($mapped['b'] ?? null) !== 20) {
    \Log::fatal('Arr::mapWithKeys 语义失败: '.var_export($mapped, true));
}

$flatAssoc = \Illuminate\Support\Arr::flatten(['App\\R' => ['default' => 'ResourceClass']]);
if (!is_array($flatAssoc) || !in_array('ResourceClass', array_values($flatAssoc), true)) {
    \Log::fatal('Arr::flatten 应展开关联嵌套: '.var_export($flatAssoc, true));
}

$u = new \Symfony\Component\String\UnicodeString('AbC');
$lower = $u->lower();
if ((string) $u !== 'AbC') {
    \Log::fatal('UnicodeString::lower 不应改写原实例');
}
if ((string) $lower !== 'abc' || $lower->length() !== 3) {
    \Log::fatal('UnicodeString::lower 语义失败');
}
$b = new \Symfony\Component\String\ByteString('hi');
if ($b->length() !== 2) {
    \Log::fatal('ByteString::length 应为字节数');
}

$uuid = \Symfony\Component\Uid\Uuid::v4();
$uuidStr = (string) $uuid;
if (!preg_match('/^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/', $uuidStr)) {
    \Log::fatal('Uuid::v4 格式失败: '.$uuidStr);
}
$ulid = \Symfony\Component\Uid\Ulid::generate();
if (strlen((string) $ulid) < 16) {
    \Log::fatal('Ulid::generate 失败');
}

$clock = new \Symfony\Component\Clock\NativeClock();
$now = $clock->now();
if (!($now instanceof \DateTimeImmutable)) {
    \Log::fatal('NativeClock::now 应返回 DateTimeImmutable, got '.gettype($now));
}

$req = \Symfony\Component\HttpFoundation\Request::create('/admin', 'GET');
$stack = new \Symfony\Component\HttpFoundation\RequestStack();
$stack->push($req);
$cur = $stack->getCurrentRequest();
if ($cur !== $req) {
    \Log::fatal('RequestStack::getCurrentRequest 失败');
}

$upload = new \Symfony\Component\HttpFoundation\File\UploadedFile(__FILE__, 'vendoraccel_smoke.php', 'text/x-php', 0, true);
if (!$upload->isValid() || $upload->getClientOriginalName() !== 'vendoraccel_smoke.php') {
    \Log::fatal('UploadedFile 语义失败');
}

$streamed = new \Symfony\Component\HttpFoundation\StreamedResponse(null, 204);
if ($streamed->getStatusCode() !== 204) {
    \Log::fatal('StreamedResponse status 失败');
}

$ir = new \Illuminate\Http\Response('hi', 200);
if ($ir->getContent() !== 'hi' || $ir->status() !== 200) {
    \Log::fatal('Illuminate\\Http\\Response 语义失败');
}

$aliases = ['name' => 'as', 'scopeBindings' => 'scope_bindings'];
$got = \Illuminate\Support\Arr::get($aliases, 'prefix', 'prefix');
if ($got !== 'prefix') {
    \Log::fatal('Arr::get 缺失键应返回 default，实际: '.var_export($got, true));
}
$gotNull = \Illuminate\Support\Arr::get(['a' => null], 'a', 'fallback');
if ($gotNull !== null) {
    \Log::fatal('Arr::get 存在的 null 值应原样返回，实际: '.var_export($gotNull, true));
}

$config = ['url' => 'mysql://x', 'driver' => 'mysql'];
$pulled = \Illuminate\Support\Arr::pull($config, 'url');
if ($pulled !== 'mysql://x') {
    \Log::fatal('Arr::pull 缺省 default 失败: '.var_export($pulled, true));
}
if (array_key_exists('url', $config)) {
    \Log::fatal('Arr::pull 未移除 url: '.var_export($config, true));
}
if (($config['driver'] ?? null) !== 'mysql') {
    \Log::fatal('Arr::pull 破坏了其余键');
}

\Log::info('vendoraccel smoke 通过');
