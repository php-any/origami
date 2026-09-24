<?php

namespace tests\php;

/**
 * Str::uuid() / Str::orderedUuid() / Str::uuid7() 必须返回 Ramsey\Uuid\UuidInterface 对象，
 * 而不是裸字符串 —— 官方调用点是 Str::orderedUuid()->toString()（telescope 每请求都会走）。
 */

use Illuminate\Support\Str;
use Ramsey\Uuid\UuidInterface;
use Ramsey\Uuid\Rfc4122\UuidInterface as Rfc4122UuidInterface;

if (!class_exists(Str::class, false)) {
    Log::info('skip: Str 原生类未注册');
    return;
}

$v4 = '/^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i';

// --- 1. 返回对象，且 instanceof 成立 -------------------------------------
$uuid = Str::uuid();
if (!is_object($uuid)) {
    Log::fatal('Str::uuid() 不是对象: ' . gettype($uuid));
}
if (!($uuid instanceof UuidInterface)) {
    Log::fatal('Str::uuid() 未实现 Ramsey\Uuid\UuidInterface');
}
if (!($uuid instanceof Rfc4122UuidInterface)) {
    Log::fatal('Str::uuid() 未实现 Ramsey\Uuid\Rfc4122\UuidInterface');
}

// --- 2. telescope 真实调用点：->toString() ------------------------------
$ordered = Str::orderedUuid()->toString();
if (!preg_match($v4, $ordered)) {
    Log::fatal('orderedUuid()->toString(): ' . $ordered);
}
if (strlen($ordered) !== 36) {
    Log::fatal('orderedUuid 长度: ' . strlen($ordered));
}

// --- 3. 字符串转换 / 拼接 / isUuid 不回归 -------------------------------
$plain = (string) Str::uuid();
if (!preg_match($v4, $plain)) {
    Log::fatal('(string) Str::uuid(): ' . $plain);
}
$cat = (string) ('id-' . Str::uuid());
if (substr($cat, 0, 3) !== 'id-' || !preg_match($v4, substr($cat, 3))) {
    Log::fatal('字符串拼接失败: ' . $cat);
}
if (!Str::isUuid((string) Str::uuid())) {
    Log::fatal('Str::isUuid((string) Str::uuid()) 为 false');
}
if (Str::isUuid('not-a-uuid')) {
    Log::fatal('Str::isUuid 误判非 uuid');
}
if (!Str::isUuid($ordered)) {
    Log::fatal('Str::isUuid(orderedUuid) 为 false');
}

// --- 4. json_encode 走 jsonSerialize -----------------------------------
$json = (string) json_encode(['id' => Str::uuid()]);
if (!preg_match('/^\{"id":"[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"\}$/', $json)) {
    Log::fatal('json_encode: ' . $json);
}
$back = json_decode($json, true);
if (!preg_match($v4, (string) $back['id'])) {
    Log::fatal('json 往返失败');
}

// --- 5. uuid7 的版本位是 7 ---------------------------------------------
$u7 = Str::uuid7()->toString();
if ($u7[14] !== '7') {
    Log::fatal('uuid7 版本位: ' . $u7);
}

// --- 6. orderedUuid 字符串按时间单调递增（时间戳前置） -----------------
$prev = '';
for ($i = 0; $i < 5; $i++) {
    $cur = (string) Str::orderedUuid();
    if ($prev !== '' && strcmp($cur, $prev) <= 0) {
        Log::fatal('orderedUuid 未单调递增: ' . $prev . ' -> ' . $cur);
    }
    $prev = $cur;
    usleep(2000);
}

// --- 7. Ramsey 实例方法面 ----------------------------------------------
if ($uuid->getHex() !== str_replace('-', '', (string) $uuid)) {
    Log::fatal('getHex: ' . $uuid->getHex());
}
if (strlen($uuid->getBytes()) !== 16) {
    Log::fatal('getBytes 长度: ' . strlen($uuid->getBytes()));
}
if ($uuid->getVersion() !== 4) {
    Log::fatal('getVersion: ' . $uuid->getVersion());
}
if ($uuid->getUrn() !== 'urn:uuid:' . (string) $uuid) {
    Log::fatal('getUrn: ' . $uuid->getUrn());
}
$other = Str::uuid();
if ($uuid->equals($other)) {
    Log::fatal('equals 误判');
}
$same = Str::uuid();
if (!$same->equals($same)) {
    Log::fatal('equals 自反失败');
}
if ($same->compareTo($same) !== 0) {
    Log::fatal('compareTo 自反失败');
}
if ($uuid->getFields()->getVersion() !== 4) {
    Log::fatal('getFields()->getVersion()');
}
// orderedUuid 仍是 v4 版本位（官方 CombGenerator 语义）
if (Str::orderedUuid()->getVersion() !== 4) {
    Log::fatal('orderedUuid 版本位不是 4');
}

// --- 8. 自定义工厂优先 --------------------------------------------------
$fixed = 'aaaaaaaa-bbbb-4ccc-8ddd-eeeeeeeeeeee';
Str::createUuidsUsing(function () use ($fixed) {
    return $fixed;
});
if ((string) Str::uuid() !== $fixed) {
    Log::fatal('createUuidsUsing 未生效: ' . Str::uuid());
}
if ((string) Str::orderedUuid() !== $fixed) {
    Log::fatal('createUuidsUsing 对 orderedUuid 未生效');
}
Str::createUuidsNormally();
if ((string) Str::uuid() === $fixed) {
    Log::fatal('createUuidsNormally 未生效');
}

// --- 9. freezeUuids 冻结住同一个 UUID -----------------------------------
$frozen = Str::freezeUuids();
if (!preg_match($v4, (string) $frozen)) {
    Log::fatal('freezeUuids 返回值: ' . $frozen);
}
if ((string) Str::uuid() !== (string) $frozen) {
    Log::fatal('freezeUuids 未冻结');
}
Str::createUuidsNormally();

// --- 10. createUuidsUsingSequence ---------------------------------------
$seq = ['11111111-1111-4111-8111-111111111111', '22222222-2222-4222-8222-222222222222'];
Str::createUuidsUsingSequence($seq);
if ((string) Str::uuid() !== $seq[0]) {
    Log::fatal('sequence[0]: ' . Str::uuid());
}
if ((string) Str::uuid() !== $seq[1]) {
    Log::fatal('sequence[1]: ' . Str::uuid());
}
if (!preg_match($v4, (string) Str::uuid())) {
    Log::fatal('sequence 用尽后未回落正常生成: ' . Str::uuid());
}
Str::createUuidsNormally();

// --- 11. 序列化载荷与上游 ramsey/uuid 一致 -------------------------------
// 上游 __serialize() 返回 ['uuid' => $this->uid]，PHP 序列化结果就是
// O:16:"Ramsey\Uuid\Uuid":1:{s:4:"uuid";s:36:"<uuid>";}
$ser = serialize(Str::uuid());
if (!preg_match('/^O:16:"Ramsey\\\\Uuid\\\\Uuid":1:\{s:4:"uuid";s:36:"[0-9a-f-]{36}";\}$/', (string) $ser)) {
    Log::fatal('serialize 载荷: ' . $ser);
}

Log::info('laravel_str_uuid_object_test ok');
